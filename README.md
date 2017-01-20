UDPX
=========
[![Build Status](https://travis-ci.org/felipejfc/udpx.svg?branch=master)](https://travis-ci.org/felipejfc/udpx)
[![Coverage Status](https://coveralls.io/repos/github/felipejfc/udpx/badge.svg)](https://coveralls.io/github/felipejfc/udpx)
[![Code Climate](https://codeclimate.com/github/felipejfc/udpx/badges/gpa.svg)](https://codeclimate.com/github/felipejfc/udpx)
[![Go Report Card](https://goreportcard.com/badge/github.com/felipejfc/udpx)](https://goreportcard.com/report/github.com/felipejfc/udpx)
[![](https://images.microbadger.com/badges/image/felipejfc/udpx.svg)](https://microbadger.com/images/felipejfc/udpx)

A Super Fast UDP Proxy that works as a NAT (has support to multiple clients) written in Golang.

### About

### Features

* Super Fast
* Can Handle Multiple Clients
* Act as a NAT
* Dynamic upstreams
* Multiple upstreams

### Dependencies
 - GO 1.7

### Compiling
```
make build
```

### Usage
```
$ ./bin/udpx --help
A fast UDP proxy that support multiple clients and dynamic upstreams

Usage:
  udpx [command]

Available Commands:
  start       starts UDP proxy
  version     Print the version number of UDPX

Use "udpx [command] --help" for more information about a command.
```

### TODO
- [x] Add config
- [x] Add command
- [x] Add tests infrastructure
- [x] Travis CI and Code Coverage
- [x] Support to multiple upstreams
- [x] Dynamically resolve upstreams
- [x] Dynamically add proxies
- [x] Dynamically remove proxies
- [x] Resolve new upstream addr if it changes
- [x] Make timeout logic faster by making less updates
- [ ] Zap has a leak, maybe use another logger
- [ ] Dynamically added proxies must be shared my multiple udpx instances
- [ ] Can persist upstreams
- [ ] Print statistics of messages sent and clients active /sec
- [ ] Persist proxy state between reboots?
- [ ] Docs
- [ ] Example
- [ ] Performance tests
- [ ] Add more tests
- [ ] Limit clients?
- [X] Docker
