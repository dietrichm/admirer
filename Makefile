.PHONY: default build test mocks

default: build test

build:
	go build ./...

test:
	go test ./...

mocks:
	go generate ./...
