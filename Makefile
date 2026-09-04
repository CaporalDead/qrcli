BIN     := qrcli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build test cover lint fmt clean install

## build: compile a static binary into ./qrcli
build:
	CGO_ENABLED=0 go build -trimpath -ldflags '$(LDFLAGS)' -o $(BIN) ./cmd/qrcli

## test: run all tests with the race detector
test:
	go test -race ./...

## cover: run tests and print per-function coverage
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## lint: run golangci-lint (see .golangci.yml)
lint:
	golangci-lint run

## fmt: format sources and tidy the module
fmt:
	gofmt -w .
	go mod tidy

## clean: remove build artifacts
clean:
	rm -f $(BIN) coverage.out

## install: install into GOBIN with the version stamped
install:
	go install -trimpath -ldflags '$(LDFLAGS)' ./cmd/qrcli
