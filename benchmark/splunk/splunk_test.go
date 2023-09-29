// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build benchmark

package splunk

import (
	"flag"
	"strings"
	"testing"

	"github.com/aws/shim-loggers-for-containerd/e2e"
	"github.com/containerd/containerd/cio"
)

const (
	splunkTokenKey              = "--splunk-token" //nolint:gosec // no real credential
	splunkURLkey                = "--splunk-url"
	splunkInsecureskipverifyKey = "--splunk-insecureskipverify"
	testSplunkURL               = "https://localhost:8089"
)

var (
	// Binary is the path the binary of the shim loggers for containerd.
	Binary      = flag.String("binary", "", "the binary of shim loggers for containerd")
	SplunkToken = flag.String("splunk-token", "", "the token to access Splunk")
)

func BenchmarkSplunk(b *testing.B) {
	testLog := strings.Repeat("a", 1024)
	args := map[string]string{
		e2e.LogDriverTypeKey:        e2e.SplunkDriverName,
		e2e.ContainerIDKey:          e2e.TestContainerID,
		e2e.ContainerNameKey:        e2e.TestContainerName,
		splunkTokenKey:              *SplunkToken,
		splunkURLkey:                testSplunkURL,
		splunkInsecureskipverifyKey: "true",
	}
	creator := cio.BinaryIO(*Binary, args)
	err := e2e.SendTestLogByContainerd(creator, testLog)
	if err != nil {
		b.Fatal(err)
	}
}
