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

package proxy_test

import (
	"time"

	. "github.com/felipejfc/udp-proxy/proxy"
	"github.com/uber-go/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Proxy", func() {

	var (
		testProxy *Proxy
	)

	BeforeEach(func() {
		debug := false
		ll := zap.ErrorLevel
		logger := zap.New(
			zap.NewJSONEncoder(),
			ll,
		)
		bindPort := 23456
		bindAddress := "localhost"
		upstreamAddr := "localhost"
		upstreamPort := 34567
		bufferSize := 4096
		connTimeout := time.Second * 1
		testProxy = GetProxy(debug, logger, bindPort, bindAddress, upstreamAddr, upstreamPort, bufferSize, connTimeout)
	})

	Describe("GetProxy", func() {
		It("Proxy should be configured", func() {
			Expect(testProxy.BindPort).To(Equal(23456))
		})
	})
})
