/*
 * Copyright (c) 2016 Felipe Cavalcanti <fjfcavalcanti@gmail.com>
 * Author: Felipe Cavalcanti <fjfcavalcanti@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package proxy

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/uber-go/zap"
)

func CheckError(err error) {
	ll := zap.ErrorLevel
	logger := zap.New(zap.NewJSONEncoder(), ll)
	if err != nil {
		logger.Fatal("error", zap.Error(err))
	}
}

type connection struct {
	udp          *net.UDPConn
	lastActivity time.Time
}

type Proxy struct {
	Logger          zap.Logger
	BindPort        int
	BindAddress     string
	UpstreamAddress string
	UpstreamPort    int
	Debug           bool
	listenerConn    *net.UDPConn
	client          *net.UDPAddr
	upstream        *net.UDPAddr
	BufferSize      int
	ConnTimeout     time.Duration
	ResolveTTL      time.Duration
	connsMap        map[string]connection
	connectionsLock *sync.RWMutex
	closed          bool
}

func GetProxy(debug bool, logger zap.Logger, bindPort int, bindAddress string, upstreamAddress string, upstreamPort int, bufferSize int, connTimeout time.Duration, resolveTTL time.Duration) *Proxy {
	proxy := &Proxy{
		Debug:           debug,
		Logger:          logger,
		BindPort:        bindPort,
		BindAddress:     bindAddress,
		BufferSize:      bufferSize,
		ConnTimeout:     connTimeout,
		UpstreamAddress: upstreamAddress,
		UpstreamPort:    upstreamPort,
		connectionsLock: new(sync.RWMutex),
		connsMap:        make(map[string]connection),
		closed:          false,
		ResolveTTL:      resolveTTL,
	}

	return proxy
}

func (p *Proxy) updateClientLastActivity(clientAddr *net.UDPAddr) {
	p.Logger.Debug("updating client last activity", zap.String("client", clientAddr.String()))
	p.connectionsLock.Lock()
	if _, found := p.connsMap[clientAddr.String()]; found {
		connWrapper := p.connsMap[clientAddr.String()]
		connWrapper.lastActivity = time.Now()
		p.connsMap[clientAddr.String()] = connWrapper
	}
	p.connectionsLock.Unlock()
}

func (p *Proxy) clientConnectionReadLoop(clientAddr *net.UDPAddr, upstreamConn *net.UDPConn) {
	for {
		buffer := make([]byte, p.BufferSize)
		size, _, err := upstreamConn.ReadFromUDP(buffer)
		if err != nil {
			p.connectionsLock.Lock()
			upstreamConn.Close()
			delete(p.connsMap, clientAddr.String())
			p.connectionsLock.Unlock()
			return
		}
		p.updateClientLastActivity(clientAddr)
		go func(data []byte, clientAddr *net.UDPAddr) {
			p.listenerConn.WriteTo(data, clientAddr)
			//p.Logger.Debug("forwarded data from upstream", zap.Int("size", size), zap.String("data", string(buffer[:size])))
		}(buffer[:size], clientAddr)
	}
}

func (p *Proxy) handlePacket(srcAddr *net.UDPAddr, data []byte, size int) {
	//p.Logger.Debug("packet received",
	//	zap.String("src address", srcAddr.String()),
	//	zap.Int("src port", srcAddr.Port),
	//	zap.String("packet", string(data[:size])),
	//)

	p.connectionsLock.RLock()
	conn, found := p.connsMap[srcAddr.String()]
	p.connectionsLock.RUnlock()

	if !found {
		conn, err := net.ListenUDP("udp", p.client)
		p.Logger.Debug("new client connection",
			zap.String("local port", conn.LocalAddr().String()),
		)

		if err != nil {
			p.Logger.Error("upd proxy failed to dial", zap.Error(err))
			return
		}

		p.connectionsLock.Lock()
		p.connsMap[srcAddr.String()] = connection{
			udp:          conn,
			lastActivity: time.Now(),
		}
		p.connectionsLock.Unlock()

		conn.WriteTo(data[:size], p.upstream)
		p.clientConnectionReadLoop(srcAddr, conn)
	} else {
		conn.udp.WriteTo(data[:size], p.upstream)
		p.connectionsLock.RLock()
		shouldUpdateLastActivity := false
		if _, found := p.connsMap[srcAddr.String()]; found {
			if p.connsMap[srcAddr.String()].lastActivity.Before(
				time.Now().Add(-p.ConnTimeout / 4)) {
				shouldUpdateLastActivity = true
			}
		}
		p.connectionsLock.RUnlock()
		if shouldUpdateLastActivity {
			p.updateClientLastActivity(srcAddr)
		}
	}
}

func (p *Proxy) readLoop() {
	for !p.closed {
		buffer := make([]byte, p.BufferSize)
		size, srcAddress, err := p.listenerConn.ReadFromUDP(buffer)
		if err != nil {
			p.Logger.Error("error", zap.Error(err))
			continue
		}
		go p.handlePacket(srcAddress, buffer, size)
	}
}

func (p *Proxy) resolveUpstreamLoop() {
	for !p.closed {
		time.Sleep(p.ResolveTTL)
		upstreamAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.UpstreamAddress, p.UpstreamPort))
		if err != nil {
			p.Logger.Error("resolve error", zap.Error(err))
			continue
		}
		if p.upstream.String() != upstreamAddr.String() {
			p.upstream = upstreamAddr
			p.Logger.Info("upstream addr changed", zap.String("upstreamAddr", p.upstream.String()))
		}
	}
}

func (p *Proxy) freeIdleSocketsLoop() {
	for !p.closed {
		time.Sleep(p.ConnTimeout)
		var clientsToTimeout []string

		p.connectionsLock.RLock()
		for client, conn := range p.connsMap {
			if conn.lastActivity.Before(time.Now().Add(-p.ConnTimeout)) {
				clientsToTimeout = append(clientsToTimeout, client)
			}
		}
		p.connectionsLock.RUnlock()

		p.connectionsLock.Lock()
		for _, client := range clientsToTimeout {
			p.Logger.Debug("client timeout", zap.String("client", client))
			p.connsMap[client].udp.Close()
			delete(p.connsMap, client)
		}
		p.connectionsLock.Unlock()
	}
}

func (p *Proxy) Close() {
	p.Logger.Warn("Closing proxy")
	p.connectionsLock.Lock()
	p.closed = true
	for _, conn := range p.connsMap {
		conn.udp.Close()
	}
	if p.listenerConn != nil {
		p.listenerConn.Close()
	}
	p.connectionsLock.Unlock()
}

func (p *Proxy) Start() {
	p.Logger.Info("Starting proxy")

	ProxyAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.BindAddress, p.BindPort))
	if err != nil {
		p.Logger.Error("error resolving bind address", zap.Error(err))
		return
	}
	p.upstream, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.UpstreamAddress, p.UpstreamPort))
	if err != nil {
		p.Logger.Error("error resolving upstream address", zap.Error(err))
	}
	p.client = &net.UDPAddr{
		IP:   ProxyAddr.IP,
		Port: 0,
		Zone: ProxyAddr.Zone,
	}
	p.listenerConn, err = net.ListenUDP("udp", ProxyAddr)
	if err != nil {
		p.Logger.Error("error listening on bind port", zap.Error(err))
		return
	}
	p.Logger.Info("UDP Proxy started!")
	if p.ConnTimeout.Nanoseconds() > 0 {
		go p.freeIdleSocketsLoop()
	} else {
		p.Logger.Warn("be warned that running without timeout to clients may be dangerous")
	}
	if p.ResolveTTL.Nanoseconds() > 0 {
		go p.resolveUpstreamLoop()
	} else {
		p.Logger.Warn("not refreshing upstream addr")
	}
	go p.readLoop()
}
