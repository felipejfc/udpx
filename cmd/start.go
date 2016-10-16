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

package cmd

import (
	"time"

	"github.com/felipejfc/udp-proxy/proxy"
	"github.com/spf13/cobra"
	"github.com/uber-go/zap"
)

var debug bool
var quiet bool

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts UDP proxy",
	Long: `Starts UDP proxy with the specified arguments. You can use
environment variables to override configuration keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		ll := zap.InfoLevel
		if debug {
			ll = zap.DebugLevel
		}
		if quiet {
			ll = zap.ErrorLevel
		}
		l := zap.New(
			zap.NewJSONEncoder(), // drop timestamps in tests
			ll,
		)

		cmdL := l.With(
			zap.Bool("debug", debug),
			zap.Bool("quiet", quiet),
		)

		cmdL.Debug("Creating proxy...")
		//TODO fix
		bindPort := 10000
		bindAddress := "0.0.0.0"
		bufferSize := 4096
		connTimeout := time.Second * 1 //millis
		upstreamAddress := "localhost"
		upstreamPort := 8830

		l1 := l.With(
			zap.String("bind address", bindAddress),
			zap.Int("bind port", bindPort),
			zap.String("upstream address", upstreamAddress),
			zap.Int("upstream port", upstreamPort),
		)

		p := proxy.GetProxy(debug, l1, bindPort, bindAddress, upstreamAddress, upstreamPort, bufferSize, connTimeout)
		cmdL.Debug("Proxy created successfully.")
		p.StartProxy()
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	startCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode (log level error)")
}
