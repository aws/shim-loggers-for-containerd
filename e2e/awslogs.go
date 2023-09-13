// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const (
	awslogsCredentialsEndpointKey = "awslogs-credentials-endpoint"
	awslogsRegionKey              = "awslogs-region"
	awslogsStreamKey              = "awslogs-stream"
	awslogsGroupKey               = "awslogs-group"
	testEcsLocalEndpointPort      = "51679"
	testAwslogsCredentialEndpoint = ":" + testEcsLocalEndpointPort + "/creds"
	testAwslogsRegion             = "us-west-2"
	testAwsLogsStream             = "test-stream"
	testAwsLogsGroup              = "test-shim-logger"
)

var testAwslogs = func() {
	// These tests are run in serial because we only define one log driver instance.
	ginkgo.Describe("awslogs shim logger", ginkgo.Serial, func() {
		ginkgo.BeforeAll(func() {
			cmd := exec.Command("docker", "run", "-d", "--name", "ecs-local-endpoint", "-p", fmt.Sprintf("%s:51679", testEcsLocalEndpointPort), "-v", "$HOME/.aws/:/home/.aws/", "-e", "AWS_REGION=us-west-2", "-e", "HOME=\"/home\"", "-e", "AWS_PROFILE=default", "-e", "ECS_LOCAL_METADATA_PORT=51679", "amazon/amazon-ecs-local-container-endpoints")
			err := cmd.Run()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.It("should send logs to awslogs log driver", func() {
			args := map[string]string{
				logDriverTypeKey:              awslogsDriverName,
				containerIdKey:                testContainerId,
				containerNameKey:              testContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwsLogsGroup,
				awslogsStreamKey:              testAwsLogsStream,
			}
			creator := cio.BinaryIO(shimLoggersBinary, args)
			// Create a new client connected to the containerd daemon
			client, err := containerd.New(containerdAddress)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			defer client.Close()
			// Create a new context with a customized namespace
			ctx := namespaces.WithNamespace(context.Background(), "testAwslogs")
			// Pull an image
			image, err := client.Pull(ctx, testImage, containerd.WithPullUnpack)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			// Create a new container with the pulled image
			container, err := client.NewContainer(ctx, testContainerId, containerd.WithImage(image),
				containerd.WithNewSnapshot("hello-snapshot", image))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			defer container.Delete(ctx, containerd.WithSnapshotCleanup)
			// Create a new task from the container and start it
			task, err := container.NewTask(ctx, creator)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			defer task.Delete(ctx)

			err = task.Start(ctx)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			statusC, err := task.Wait(ctx)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			// Run the "echo" command in the container to print "test-e2e-log" to the log
			_, err = task.Exec(ctx, testExecId, &specs.Process{
				Args: []string{"/bin/sh", "-c", "echo 'test-e2e-log'"},
			}, creator)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			// Waiting for the task to finish
			status := <-statusC
			code, _, err := status.Result()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(code).Should(gomega.Equal(0))
		})
	})
}
