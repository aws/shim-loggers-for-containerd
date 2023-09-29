// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build e2e

// Package e2e provides e2e tests of shim loggers for containerd.
package e2e

import (
	"context"
	"errors"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

const (
	// LogDriver options.

	// LogDriverTypeKey is the key of log driver type.
	LogDriverTypeKey = "--log-driver"
	// AwslogsDriverName is the name of awslogs driver.
	AwslogsDriverName = "awslogs"
	// FluentdDriverName is the name of fluentd driver.
	FluentdDriverName = "fluentd"
	// SplunkDriverName is the name of splunk driver.
	SplunkDriverName = "splunk"
	// ContainerIDKey is the key of the container id.
	ContainerIDKey = "--container-id"
	// ContainerNameKey is the key of the container name.
	ContainerNameKey = "--container-name"
	// TestContainerID is the id of the tes container.
	TestContainerID = "210987654321"
	// TestContainerName is the name of the test container.
	TestContainerName                = "test-container-name"
	containerdAddress                = "/run/containerd/containerd.sock"
	testImage                        = "public.ecr.aws/docker/library/ubuntu:latest"
	testLogPrefix                    = "test-e2e-log-"
	containerdTaskExitNonZeroMessage = "\"containerd task exits with non-zero\""
)

// SendTestLogByContainerd sends a testLog to a specific shim logger by containerd.
func SendTestLogByContainerd(creator cio.Creator, testLog string) error {
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
	container, err := client.NewContainer(ctx, TestContainerID, containerd.WithImage(image),
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
