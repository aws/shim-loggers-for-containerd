// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	shimLoggersBinary = "/Users/ningziwe/shim-loggers-for-containerd/shim-loggers-for-containerd"
	// LogDriver options
	logDriverTypeKey  = "log-driver"
	awslogsDriverName = "awslogs"
	fluentdDriverName = "fluentd"
	splunkDriverName  = "splunk"
	containerIdKey    = "container-id"
	containerNameKey  = "container-name"
	testContainerId   = "test-container-id"
	testContainerName = "test-container-name"
	containerdAddress = "/run/containerd/containerd.sock"
	testImage         = "public.ecr.aws/docker/library/ubuntu:latest"
	testExecId        = "test-exec-id"
)

func TestShimLoggers(t *testing.T) {
	const description = "Shim loggers for containerd E2E Tests"

	ginkgo.Describe("", func() {
		testAwslogs()
		testFluentd()
		testSplunk()
	})

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, description)
}
