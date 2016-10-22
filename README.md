UDPX
=========
[![Build Status](https://travis-ci.org/felipejfc/udpx.svg?branch=master)](https://travis-ci.org/felipejfc/udpx)
[![Coverage Status](https://coveralls.io/repos/github/felipejfc/udpx/badge.svg)](https://coveralls.io/github/felipejfc/udpx)
[![Code Climate](https://codeclimate.com/github/felipejfc/udpx/badges/gpa.svg)](https://codeclimate.com/github/felipejfc/udpx)
[![Go Report Card](https://goreportcard.com/badge/github.com/felipejfc/udpx)](https://goreportcard.com/report/github.com/felipejfc/udpx)

A Super Fast UDP Proxy that works as a NAT (has support to multiple clients) written in Golang.

### Features

* Super Fast
* Can Handle Multiple Clients
* Act as a NAT
* Dynamic upstreams
* Multiple upstreams

### TODO
- [x] Add config
- [x] Add command
- [x] Add tests infrastructure
- [x] Travis CI and Code Coverage
- [x] Support to multiple upstreams
- [ ] Dynamically resolve upstreams
- [x] Dynamically add proxies
- [x] Dynamically remove proxies
- [x] Resolve new upstream addr if it changes
- [ ] Dynamically added proxies must be shared my multiple udpx instances
- [ ] Can persist upstreams
- [ ] Persist proxy state between reboots?
- [ ] Docs
- [ ] Performance tests
- [ ] Add more tests
- [X] Docker
