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
	"os"

	"github.com/uber-go/zap"
)

func CheckError(err error) {
	ll := zap.ErrorLevel
	logger := zap.New(zap.NewJSONEncoder(), ll)
	if err != nil {
		logger.Error("error", zap.Error(err))
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

type Proxy struct {
	Logger        zap.Logger
	LocalPort     int
	BindAddress   string
	Debug         bool
	ProxyConn     *net.UDPConn
	BufferSize    int
	SocketTimeout int
}

func GetProxy(debug bool, logger zap.Logger, localPort int, bindAddress string, bufferSize int, socketTimeout int) *Proxy {
	proxy := &Proxy{
		Debug:         debug,
		Logger:        logger,
		LocalPort:     localPort,
		BindAddress:   bindAddress,
		BufferSize:    bufferSize,
		SocketTimeout: socketTimeout,
	}

	return proxy
}

func (p *Proxy) readLoop() {
	for {
		buffer := make([]byte, p.BufferSize)
		size, srcAddress, err := p.ProxyConn.ReadFromUDP(buffer)
		if err != nil {
			//TODO Should I only log the error here?
			p.Logger.Error("error", zap.Error(err))
		}
		p.Logger.Debug("packet received",
			zap.String("src address", srcAddress.String()),
			zap.Int("src port", srcAddress.Port),
			zap.String("packet", string(buffer[:size])),
		)
	}
}

func (p *Proxy) freeSocketsLoop() {

}

func (p *Proxy) StartProxy() {
	p.Logger.Info("Starting proxy",
		zap.String("BindAddress", p.BindAddress),
		zap.Int("LocalPort", p.LocalPort),
	)

	ProxyAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.BindAddress, p.LocalPort))
	CheckError(err)
	p.ProxyConn, err = net.ListenUDP("udp", ProxyAddr)
	CheckError(err)
	p.Logger.Info("UDP Proxy listening...")
	p.readLoop()
}
