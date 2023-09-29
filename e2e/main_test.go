// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build e2e

package e2e

import (
	"flag"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var (
	// Binary is the path the binary of the shim loggers for containerd.
	Binary      = flag.String("binary", "", "the binary of shim loggers for containerd")
	LogDriver   = flag.String("log-driver", "", "the log driver to test")
	SplunkToken = flag.String("splunk-token", "", "the token to access Splunk")
)

func TestShimLoggers(t *testing.T) {
	t.Parallel()

	const description = "Shim loggers for containerd E2E Tests"

	ginkgo.Describe("", func() {
		if *LogDriver == AwslogsDriverName || *LogDriver == "" {
			testAwslogs()
		}
		if *LogDriver == FluentdDriverName || *LogDriver == "" {
			testFluentd()
		}
		if *LogDriver == SplunkDriverName || *LogDriver == "" {
			testSplunk(*SplunkToken)
		}
	})

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, description)
}
