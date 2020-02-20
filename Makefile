ROOT := $(shell pwd)

SOURCEDIR=./

AWS_CONTAINERD_LOGGERS_DIR=$(SOURCEDIR)
AWS_CONTAINERD_LOGGERS_BINARY=$(ROOT)/bin/shim-loggers-for-containerd
SOURCES=$(shell find $(SOURCEDIR) -name '*.go')

BINPATH:=$(abspath ./bin)

.PHONY: all
all: build test

.PHONY: build
build: $(AWS_CONTAINERD_LOGGERS_BINARY)

$(AWS_CONTAINERD_LOGGERS_BINARY):
	go build -o $(AWS_CONTAINERD_LOGGERS_BINARY) $(AWS_CONTAINERD_LOGGERS_DIR)

test: $(SOURCES)
	go test -tags unit -race -timeout 30s -cover $(shell go list ./...) --count=1

.PHONY: clean
clean:
	@rm -f $(BINPATH)/*
