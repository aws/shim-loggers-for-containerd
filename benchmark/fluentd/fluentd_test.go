// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build benchmark

package fluentd

import (
	"flag"
	"strings"
	"testing"

	"github.com/aws/shim-loggers-for-containerd/e2e"
	"github.com/containerd/containerd/cio"
)

var (
	// Binary is the path the binary of the shim loggers for containerd.
	Binary = flag.String("binary", "", "the binary of shim loggers for containerd")
)

func BenchmarkFluentd(b *testing.B) {
	testLog := strings.Repeat("a", 1024)
	args := map[string]string{
		e2e.LogDriverTypeKey: e2e.FluentdDriverName,
		e2e.ContainerIDKey:   e2e.TestContainerID,
		e2e.ContainerNameKey: e2e.TestContainerName,
	}
	creator := cio.BinaryIO(*Binary, args)
	err := e2e.SendTestLogByContainerd(creator, testLog)
	if err != nil {
		b.Fatal(err)
	}
}
