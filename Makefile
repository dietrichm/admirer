.PHONY: default build test

default: build test

build:
	go build ./...

test:
	go test ./...
