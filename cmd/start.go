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
	"os"

	"github.com/felipejfc/udpx/api"
	"github.com/felipejfc/udpx/proxy"
	"github.com/spf13/cobra"
	"github.com/uber-go/zap"
)

var debug bool
var quiet bool
var useAPI bool
var bindAddress string
var bufferSize int
var apiBindPort int
var configPath string

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

		cmdL.Debug("Creating proxies...")

		proxyConfigs := proxy.LoadProxyConfigsFromConfigFiles(configPath)

		if len(proxyConfigs) == 0 {
			if !useAPI {
				cmdL.Fatal("no proxy config loaded")
			} else {
				cmdL.Warn("no proxy config loaded")
			}
		}

		pm := proxy.GetManager()
		pm.Configure(debug, l, bindAddress, bufferSize)

		for _, proxyConfig := range proxyConfigs {
			//TODO guardar proxies e verificar conflitos de bind port
			if pm.RegisterProxy(proxyConfig) != true {
				cmdL.Warn("proxy already loaded with the same bind port", zap.Int("bindPort", proxyConfig.BindPort))
			}
		}

		if useAPI {
			ll := l.With(
				zap.String("bind address", bindAddress),
				zap.Int("bind port", apiBindPort),
			)
			a := api.GetAPI(bindAddress, apiBindPort, debug, ll)
			a.Start()
		}

		exitSignal := make(chan os.Signal)
		<-exitSignal
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVarP(&bufferSize, "bufferSize", "B", 4096, "Datagrams buffer size")
	startCmd.Flags().IntVarP(&apiBindPort, "apiBindPort", "p", 8080, "The port that udpx api will bind to")
	startCmd.Flags().StringVarP(&bindAddress, "bind", "b", "0.0.0.0", "Host to bind proxies and api")
	startCmd.Flags().StringVarP(&configPath, "configPath", "c", "./config", "Path to the folder containing the config files")
	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	startCmd.Flags().BoolVarP(&useAPI, "api", "a", false, "Start udpx api for managing upstreams dinamically")
	startCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode (log level error)")
}
