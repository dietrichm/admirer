.PHONY: default test build mocks vendor

default: test build

test:
	go test ./...

build:
	go build ./...

mocks:
	go generate ./...

vendor:
	go get
	go mod vendor
