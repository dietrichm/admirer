.PHONY: default build test mocks vendor

default: build test

build:
	go build ./...

test:
	go test ./...

mocks:
	go generate ./...

vendor:
	go get
	go mod vendor
