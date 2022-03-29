/*
 * Copyright (c) 2020 Felipe Cavalcanti <fjfcavalcanti@gmail.com>
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
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

// CheckError checks for error
func CheckError(err error) {
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("error", zap.Error(err))
	}
}

type connection struct {
	udp          *net.UDPConn
	lastActivity time.Time
}

type packet struct {
	src     *net.UDPAddr
	data    []byte
	dataLen int
}

// Proxy struct
type Proxy struct {
	Logger                 *zap.Logger
	BindPort               int
	BindAddress            string
	UpstreamAddress        string
	UpstreamPort           int
	Debug                  bool
	listenerConn           *net.UDPConn
	client                 *net.UDPAddr
	upstream               *net.UDPAddr
	BufferSize             int
	ConnTimeout            time.Duration
	ResolveTTL             time.Duration
	connsMap               sync.Map
	closed                 bool
	clientMessageChannel   chan (packet)
	upstreamMessageChannel chan (packet)
	bufferPool             sync.Pool
}

// GetProxy gets the proxy
func GetProxy(debug bool, logger *zap.Logger, bindPort int, bindAddress string, upstreamAddress string, upstreamPort int, bufferSize int, connTimeout time.Duration, resolveTTL time.Duration) *Proxy {
	proxy := &Proxy{
		Debug:                  debug,
		Logger:                 logger,
		BindPort:               bindPort,
		BindAddress:            bindAddress,
		BufferSize:             bufferSize,
		ConnTimeout:            connTimeout,
		UpstreamAddress:        upstreamAddress,
		UpstreamPort:           upstreamPort,
		closed:                 false,
		ResolveTTL:             resolveTTL,
		clientMessageChannel:   make(chan packet),
		upstreamMessageChannel: make(chan packet),
		bufferPool:             sync.Pool{New: func() interface{} { return make([]byte, bufferSize) }},
	}

	return proxy
}

func (p *Proxy) updateClientLastActivity(clientAddrString string) {
	p.Logger.Debug("updating client last activity", zap.String("client", clientAddrString))
	if connWrapper, found := p.connsMap.Load(clientAddrString); found {
		connWrapper.(*connection).lastActivity = time.Now()
	}
}

func (p *Proxy) clientConnectionReadLoop(clientAddr *net.UDPAddr, upstreamConn *net.UDPConn) {
	clientAddrString := clientAddr.String()
	for {
		msg := p.bufferPool.Get().([]byte)
		size, _, err := upstreamConn.ReadFromUDP(msg[0:])
		if err != nil {
			upstreamConn.Close()
			p.connsMap.Delete(clientAddrString)
			return
		}
		p.updateClientLastActivity(clientAddrString)
		p.upstreamMessageChannel <- packet{
			src:     clientAddr,
			data:    msg,
			dataLen: size,
		}
	}
}

func (p *Proxy) handlerUpstreamPackets() {
	for pa := range p.upstreamMessageChannel {
		p.Logger.Debug("forwarded data from upstream", zap.Int("size", pa.dataLen), zap.String("data", string(pa.data[:pa.dataLen])))
		p.listenerConn.WriteTo(pa.data[:pa.dataLen], pa.src)
		p.bufferPool.Put(pa.data)
	}
}

func (p *Proxy) handleClientPackets() {
	for pa := range p.clientMessageChannel {
		packetSourceString := pa.src.String()
		p.Logger.Debug("packet received",
			zap.String("src address", packetSourceString),
			zap.Int("src port", pa.src.Port),
			zap.String("packet", string(pa.data[:pa.dataLen])),
			zap.Int("size", pa.dataLen),
		)

		conn, found := p.connsMap.Load(packetSourceString)
		if !found {
			conn, err := net.ListenUDP("udp", p.client)
			p.Logger.Debug("new client connection",
				zap.String("local port", conn.LocalAddr().String()),
			)

			if err != nil {
				p.Logger.Error("upd proxy failed to dial", zap.Error(err))
				return
			}

			p.connsMap.Store(packetSourceString, &connection{
				udp:          conn,
				lastActivity: time.Now(),
			})

			conn.WriteTo(pa.data[:pa.dataLen], p.upstream)
			go p.clientConnectionReadLoop(pa.src, conn)
		} else {
			conn.(*connection).udp.WriteTo(pa.data[:pa.dataLen], p.upstream)
			shouldUpdateLastActivity := false
			if conn, found := p.connsMap.Load(packetSourceString); found {
				if conn.(*connection).lastActivity.Before(
					time.Now().Add(-p.ConnTimeout / 4)) {
					shouldUpdateLastActivity = true
				}
			}
			if shouldUpdateLastActivity {
				p.updateClientLastActivity(packetSourceString)
			}
		}
		p.bufferPool.Put(pa.data)
	}
}

func (p *Proxy) readLoop() {
	for !p.closed {
		msg := p.bufferPool.Get().([]byte)
		size, srcAddress, err := p.listenerConn.ReadFromUDP(msg[0:])
		if err != nil {
			p.Logger.Error("error", zap.Error(err))
			continue
		}
		p.clientMessageChannel <- packet{
			src:     srcAddress,
			data:    msg,
			dataLen: size,
		}
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

		p.connsMap.Range(func(k, conn interface{}) bool {
			if conn.(*connection).lastActivity.Before(time.Now().Add(-p.ConnTimeout)) {
				clientsToTimeout = append(clientsToTimeout, k.(string))
			}
			return true
		})

		for _, client := range clientsToTimeout {
			p.Logger.Debug("client timeout", zap.String("client", client))
			conn, ok := p.connsMap.Load(client)
			if ok {
				conn.(*connection).udp.Close()
				p.connsMap.Delete(client)
			}
		}
	}
}

// Close stops the proxy
func (p *Proxy) Close() {
	p.Logger.Warn("Closing proxy")
	p.closed = true
	p.connsMap.Range(func(k, conn interface{}) bool {
		conn.(*connection).udp.Close()
		return true
	})
	if p.listenerConn != nil {
		p.listenerConn.Close()
	}
}

// Start starts the proxy
func (p *Proxy) Start() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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
	for i := 0; i < runtime.NumCPU(); i++ {
		go p.readLoop()
		go p.handleClientPackets()
		go p.handlerUpstreamPackets()
	}
}
