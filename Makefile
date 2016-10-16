.PHONY: all

build:
	@go build -o udp-proxy.o

run:
	@go run main.go start

test:
	@ginkgo -r --cover .
