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

package api

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/uber-go/zap"
)

type Api struct {
	BindAddress string
	BindPort    int
	http        *echo.Echo
	logger      zap.Logger
	debug       bool
}

func GetApi(bindAddress string, bindPort int, debug bool, logger zap.Logger) *Api{
	api := &Api{
		BindAddress: bindAddress,
		BindPort:    bindPort,
		http:        echo.New(),
		debug:       debug,
		logger:      logger,
	}
	api.configureApi()
	return api
}

func (a *Api) configureApi() {
	a.http.GET("/healthcheck", HealthCheckHandler)
	a.logger.Debug("api configured!")
}

func (a *Api) Start() {
	go a.http.Run(standard.New(fmt.Sprintf(":%d", a.BindPort)))
	a.logger.Info("api started!")
}
