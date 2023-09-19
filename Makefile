# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ROOT := $(shell pwd)

SOURCEDIR=./

exesuffix=$(go env GOEXE)

AWS_CONTAINERD_LOGGERS_DIR=$(SOURCEDIR)
AWS_CONTAINERD_LOGGERS_BINARY=$(ROOT)/bin/shim-loggers-for-containerd$(exesuffix)
SOURCES=$(shell find $(SOURCEDIR) -name '*.go')

BINPATH:=$(abspath ./bin)
DEPSPATH:=$(abspath ./deps)

.PHONY: all
all: build test

.PHONY: build
build: $(AWS_CONTAINERD_LOGGERS_BINARY)

$(AWS_CONTAINERD_LOGGERS_BINARY):
	go build -o $(AWS_CONTAINERD_LOGGERS_BINARY) $(AWS_CONTAINERD_LOGGERS_DIR)

.PHONY: test-unit
test-unit: $(SOURCES)
	go test -tags unit -race -timeout 30s -cover $(shell go list ./... | grep -v e2e) --count=1

.PHONY: test-e2e
test-e2e:
	go test -timeout 30m ./e2e -test.v -ginkgo.v --binary "$(AWS_CONTAINERD_LOGGERS_BINARY)"

.PHONY: test-e2e-for-awslogs
test-e2e-for-aws-logs:
	go test -timeout 30m ./e2e -test.v -ginkgo.v --binary "$(AWS_CONTAINERD_LOGGERS_BINARY)" --log-driver "awslogs"

.PHONY: test-e2e-for-fluentd
test-e2e-for-fluentd:
	go test -timeout 30m ./e2e -test.v -ginkgo.v --binary "$(AWS_CONTAINERD_LOGGERS_BINARY)" --log-driver "fluentd"

.PHONY: test-e2e-for-splunk
test-e2e-for-splunk:
	go test -timeout 30m ./e2e -test.v -ginkgo.v --binary "$(AWS_CONTAINERD_LOGGERS_BINARY)" --log-driver "splunk" --splunk-token ${SPLUNK_TOKEN}

.PHONY: coverage
coverage:
	go test -tags unit $(shell go list ./... | grep -v e2e) -coverprofile=test-coverage.out
	go tool cover -html=test-coverage.out

.PHONY: lint
lint: $(SOURCES)
	$(DEPSPATH)/golangci-lint run

.PHONY: mdlint
# Install it locally: https://github.com/igorshubovych/markdownlint-cli#installation
# Or see `mdlint-ctr` below or https://github.com/DavidAnson/markdownlint#related.
mdlint:
	markdownlint '**/*.md'

.PHONY: mdlint-ctr
# If markdownlint is not installed, you can run markdownlint within a container.
mdlint-ctr:
	docker run --rm -v "$(shell pwd):/repo:ro" -w /repo avtodev/markdown-lint:v1 '**/*.md'

.get-deps-stamp:
	GO111MODULE=off GOBIN=$(DEPSPATH) go get golang.org/x/tools/cmd/goimports
	GOBIN=$(DEPSPATH) go get github.com/golang/mock/mockgen@v1.4.1
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(DEPSPATH) v1.54.0
	$(DEPSPATH)/golangci-lint --version
	touch .get-deps-stamp

.PHONY: get-deps
get-deps: .get-deps-stamp

.PHONY: clean
clean:
	@rm -f $(BINPATH)/*
	@rm -f .get-deps-stamp
