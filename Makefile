.PHONY: default check test build mocks vendor

default: check test build

check:
	go vet ./...

test:
	go test ./...

build:
	go build ./...

mocks:
	go generate ./...

vendor:
	go get
	go mod vendor
