#Copyright (c) 2016 Felipe Cavalcanti <fjfcavalcanti@gmail.com>
#Author: Felipe Cavalcanti <fjfcavalcanti@gmail.com>
#
#Permission is hereby granted, free of charge, to any person obtaining a copy of
#this software and associated documentation files (the "Software"), to deal in
#the Software without restriction, including without limitation the rights to
#use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
#the Software, and to permit persons to whom the Software is furnished to do so,
#subject to the following conditions:
#
#The above copyright notice and this permission notice shall be included in all
#copies or substantial portions of the Software.
#
#THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
#FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
#COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
#IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
#CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

FROM golang:1.7.3-alpine

MAINTAINER Felipe Cavalcanti <fjfcavalcanti@gmail.com>

RUN apk update
RUN apk add make

RUN mkdir -p /go/src/github.com/felipejfc/udpx
WORKDIR /go/src/github.com/felipejfc/udpx

ADD . /go/src/github.com/felipejfc/udpx
RUN make build-cross-linux

RUN mkdir /app
RUN mv /go/src/github.com/felipejfc/udpx/bin/udpx-linux-x86_64 /app/udpx
RUN mv /go/src/github.com/felipejfc/udpx/config /app/config
RUN rm -r /go/src/github.com/felipejfc/udpx

WORKDIR /app

EXPOSE 8080
VOLUME /app/config

CMD /app/udpx start -c /app/config -a
