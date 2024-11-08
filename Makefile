all: lint test build

build:
	go build --mod=vendor -o microbatcher

test:
	go test -v --mod=vendor -cover -coverprofile cover.out ./...

cov-html: test
	go tool cover -html=cover.out -o cover.html

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

vendor:
	go mod tidy
	go mod vendor

.PHONY: all build test lint vendor cov-html lint-fix
