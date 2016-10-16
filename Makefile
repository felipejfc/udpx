# Copyright (c) 2016 Felipe Cavalcanti <fjfcavalcanti@gmail.com>
# Author: Felipe Cavalcanti <fjfcavalcanti@gmail.com>
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

.PHONY: all

build:
	@mkdir -p bin
	@go build -o bin/udpx

run:
	@go run main.go start

test:
	@ginkgo -r --cover .

coverage: test
	@echo "mode: count" > coverage-all.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> coverage-all.out; done'

coverage-html: coverage
	@go tool cover -html=coverage-all.out

setup-ci:
	@go get -u github.com/Masterminds/glide/...
	@go get github.com/mattn/goveralls
	@go get github.com/onsi/ginkgo/ginkgo
	@glide install
