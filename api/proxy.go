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
	"net/http"

	"github.com/felipejfc/udpx/proxy"
	"github.com/labstack/echo"
)

func NewProxyHandler(c echo.Context) error {
	p := new(proxy.ProxyInstance)
	pm := proxy.GetManager()
	if err := c.Bind(p); err != nil {
		return err
	}
	if p.BindPort == 0 {
		return c.String(http.StatusUnprocessableEntity, "bindPort required")
	}
	if p.UpstreamPort == 0 {
		return c.String(http.StatusUnprocessableEntity, "upstreamPort required")
	}
	if p.UpstreamAddress == "" {
		return c.String(http.StatusUnprocessableEntity, "upstreamAddress required")
	}
	if p.Name == "" {
		return c.String(http.StatusUnprocessableEntity, "name required")
	}
	if pm.RegisterProxy(*p) != true {
		return c.String(http.StatusConflict, fmt.Sprintf("some proxy might already be listening on port %d", p.BindPort))
	}
	return c.JSON(http.StatusCreated, p)
}

func GetProxyByBindPortHandler(c echo.Context) error {
	pm := proxy.GetManager()
	p := pm.GetConfigByBindPort(c.Param("port"))
	if p == nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, p)
}

func UnregisterProxyByPortHandler(c echo.Context) error {
	pm := proxy.GetManager()
	success := pm.UnregisterByBindPort(c.Param("port"))
	if success == false {
		return echo.ErrNotFound
	}
	return c.String(http.StatusOK, "OK")
}
