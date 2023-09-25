// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	// LogDriver options.
	logDriverTypeKey                 = "--log-driver"
	awslogsDriverName                = "awslogs"
	fluentdDriverName                = "fluentd"
	splunkDriverName                 = "splunk"
	containerIDKey                   = "--container-id"
	containerNameKey                 = "--container-name"
	testContainerID                  = "test-container-id"
	testContainerName                = "test-container-name"
	containerdAddress                = "/run/containerd/containerd.sock"
	testImage                        = "public.ecr.aws/docker/library/ubuntu:latest"
	testLogPrefix                    = "test-e2e-log-"
	containerdTaskExitNonZeroMessage = "\"containerd task exits with non-zero\""
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
		if *LogDriver == awslogsDriverName || *LogDriver == "" {
			testAwslogs()
		}
		if *LogDriver == fluentdDriverName || *LogDriver == "" {
			testFluentd()
		}
		if *LogDriver == splunkDriverName || *LogDriver == "" {
			testSplunk(*SplunkToken)
		}
	})

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, description)
}

func sendTestLogByContainerd(creator cio.Creator, testLog string) error {
	// Create a new client connected to the containerd daemon
	client, err := containerd.New(containerdAddress)
	if err != nil {
		return err
	}
	defer client.Close() //nolint:errcheck // closing client
	// Create a new context with a customized namespace
	ctx := namespaces.WithNamespace(context.Background(), "testShimLoggers")
	// Pull an image
	image, err := client.Pull(ctx, testImage, containerd.WithPullUnpack)
	if err != nil {
		return err
	} // Create a new container with the pulled image
	container, err := client.NewContainer(ctx, testContainerID, containerd.WithImage(image),
		containerd.WithNewSnapshot("test-snapshot", image), containerd.WithNewSpec(oci.WithImageConfig(image),
			oci.WithProcessArgs("/bin/sh", "-c", fmt.Sprintf("printf \"%s\"", testLog))))
	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup) //nolint:errcheck // testing only
	// Create a new task from the container and start it
	task, err := container.NewTask(ctx, creator)
	if err != nil {
		return err
	}
	defer task.Delete(ctx) //nolint:errcheck // testing only

	err = task.Start(ctx)
	if err != nil {
		return err
	}

	statusC, err := task.Wait(ctx)
	if err != nil {
		return err
	}
	// Waiting for the task to finish
	status := <-statusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	if code != uint32(0) {
		return errors.New(containerdTaskExitNonZeroMessage)
	}
	return nil
}
