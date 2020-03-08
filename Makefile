ROOT := $(shell pwd)

SOURCEDIR=./

AWS_CONTAINERD_LOGGERS_DIR=$(SOURCEDIR)
AWS_CONTAINERD_LOGGERS_BINARY=$(ROOT)/bin/shim-loggers-for-containerd
SOURCES=$(shell find $(SOURCEDIR) -name '*.go')

BINPATH:=$(abspath ./bin)
DEPSPATH:=$(abspath ./deps)

.PHONY: all
all: build test

.PHONY: build
build: $(AWS_CONTAINERD_LOGGERS_BINARY)

$(AWS_CONTAINERD_LOGGERS_BINARY):
	go build -o $(AWS_CONTAINERD_LOGGERS_BINARY) $(AWS_CONTAINERD_LOGGERS_DIR)

test: $(SOURCES)
	go test -tags unit -race -timeout 30s -cover $(shell go list ./...) --count=1

.PHONY: lint
lint: $(SOURCES)
	$(DEPSPATH)/golangci-lint run

.get-deps-stamp:
	GO111MODULE=off GOBIN=$(DEPSPATH) go get golang.org/x/tools/cmd/goimports
	GOBIN=$(DEPSPATH) go get github.com/golang/mock/mockgen@v1.3.1
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(DEPSPATH) v1.21.0
	$(DEPSPATH)/golangci-lint --version
	touch .get-deps-stamp

.PHONY: get-deps
get-deps: .get-deps-stamp

.PHONY: clean
clean:
	@rm -f $(BINPATH)/*
	@rm -f .get-deps-stamp
