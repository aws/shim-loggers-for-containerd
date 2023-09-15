// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"flag"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	// LogDriver options
	logDriverTypeKey  = "--log-driver"
	awslogsDriverName = "awslogs"
	containerIdKey    = "--container-id"
	containerNameKey  = "--container-name"
	testContainerId   = "test-container-id"
	testContainerName = "test-container-name"
	containerdAddress = "/run/containerd/containerd.sock"
	testImage         = "public.ecr.aws/docker/library/ubuntu:latest"
)

var (
	// Binary is the path the binary of the shim loggers for containerd
	Binary    = flag.String("binary", "", "the binary of shim loggers for containerd")
	LogDriver = flag.String("log-driver", "", "the log driver to test")
)

func TestShimLoggers(t *testing.T) {
	const description = "Shim loggers for containerd E2E Tests"

	ginkgo.Describe("", func() {
		if *LogDriver == awslogsDriverName || *LogDriver == "" {
			testAwslogs()
		}
	})

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, description)
}
