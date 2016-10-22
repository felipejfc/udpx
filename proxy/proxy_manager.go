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
	"strconv"
	"sync"
	"time"

	"github.com/uber-go/zap"
)

type Manager struct {
	Debug                bool
	Logger               zap.Logger
	BindAddress          string
	BufferSize           int
	DefaultClientTimeout int
	DefaultResolveTTL    int
}

var ProxyConfigStorage = make(map[string]*ProxyInstance)
var ProxyStorage = make(map[string]*Proxy)
var instance *Manager
var once sync.Once

func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{}
	})
	return instance
}

func (p *Manager) Configure(debug bool, logger zap.Logger, bindAddress string, bufferSize int, defaultClientTimeout int, defaultResolveTTL int) {
	p.Logger = logger
	p.Debug = debug
	p.BindAddress = bindAddress
	p.BufferSize = bufferSize
	p.DefaultClientTimeout = defaultClientTimeout
	p.DefaultResolveTTL = defaultResolveTTL
	p.Logger.Info("proxy manager configured!", zap.Bool("debug", p.Debug), zap.String("bindAddress", p.BindAddress), zap.Int("bufferSize", p.BufferSize), zap.Int("defaultResolveTTL", defaultResolveTTL), zap.Int("defaultClientTimeout", defaultClientTimeout))
}

func (p *Manager) RegisterProxy(proxyInstance ProxyInstance) bool {
	bindPortString := strconv.Itoa(proxyInstance.BindPort)
	_, found := ProxyConfigStorage[bindPortString]
	if found {
		return false
	}
	ProxyConfigStorage[bindPortString] = &proxyInstance
	if proxyInstance.ClientTimeout == 0 {
		proxyInstance.ClientTimeout = p.DefaultClientTimeout
	}
	if proxyInstance.ResolveTTL == 0 {
		proxyInstance.ResolveTTL = p.DefaultResolveTTL
	}
	ll := p.Logger.With(
		zap.String("bind address", p.BindAddress),
		zap.Int("bind port", proxyInstance.BindPort),
		zap.String("upstream address", proxyInstance.UpstreamAddress),
		zap.Int("upstream port", proxyInstance.UpstreamPort),
		zap.String("name", proxyInstance.Name),
		zap.Int("resolveTTL", proxyInstance.ResolveTTL),
		zap.Int("clientTimeout", proxyInstance.ClientTimeout),
	)
	pp := GetProxy(p.Debug, ll, proxyInstance.BindPort, p.BindAddress, proxyInstance.UpstreamAddress, proxyInstance.UpstreamPort, p.BufferSize, time.Duration(proxyInstance.ClientTimeout)*time.Millisecond, time.Duration(proxyInstance.ResolveTTL)*time.Millisecond)
	ProxyStorage[bindPortString] = pp
	pp.Start()
	return true
}

func (p *Manager) GetConfigByBindPort(port string) *ProxyInstance {
	pi, _ := ProxyConfigStorage[port]
	return pi
}

func (p *Manager) UnregisterByBindPort(port string) bool {
	pp, _ := ProxyStorage[port]
	if pp == nil {
		return false
	}
	pp.Close()
	delete(ProxyStorage, port)
	delete(ProxyConfigStorage, port)
	return true
}

func (p *Manager) PersistProxyConfig(proxy *ProxyInstance) error {
	//TODO
	return nil
}
