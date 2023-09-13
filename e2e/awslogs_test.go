// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	awslogsCredentialsEndpointKey = "--awslogs-credentials-endpoint"
	awslogsRegionKey              = "--awslogs-region"
	awslogsStreamKey              = "--awslogs-stream"
	awslogsGroupKey               = "--awslogs-group"
	testEcsLocalEndpointPort      = "51679"
	testAwslogsCredentialEndpoint = ":" + testEcsLocalEndpointPort + "/creds"
	testAwslogsRegion             = "us-west-2"
	testAwsLogsStream             = "test-stream"
	testAwsLogsGroup              = "test-shim-logger"
	testAwsLogsMessage            = "test-e2e-log"
)

var testAwslogs = func() {
	// These tests are run in serial because we only define one log driver instance.
	ginkgo.Describe("awslogs shim logger", ginkgo.Serial, func() { //nolint:typecheck
		var cwClient *cloudwatchlogs.Client
		ginkgo.BeforeEach(func() {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(testAwslogsRegion))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			cwClient = cloudwatchlogs.NewFromConfig(cfg)
			_, err = cwClient.DeleteLogStream(context.TODO(), &cloudwatchlogs.DeleteLogStreamInput{
				LogStreamName: aws.String(testAwsLogsStream),
				LogGroupName:  aws.String(testAwsLogsGroup),
			})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.AfterEach(func() {
			_, err := cwClient.DeleteLogStream(context.TODO(), &cloudwatchlogs.DeleteLogStreamInput{
				LogStreamName: aws.String(testAwsLogsStream),
				LogGroupName:  aws.String(testAwsLogsGroup),
			})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.It("should send logs to awslogs log driver", func() { //nolint:typecheck
			args := map[string]string{
				logDriverTypeKey:              awslogsDriverName,
				containerIdKey:                testContainerId,
				containerNameKey:              testContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwsLogsGroup,
				awslogsStreamKey:              testAwsLogsStream,
			}
			creator := cio.BinaryIO(*Binary, args)
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
				containerd.WithNewSnapshot("test-snapshot", image), containerd.WithNewSpec(oci.WithImageConfig(image),
					oci.WithProcessArgs("/bin/sh", "-c", fmt.Sprintf("echo '%s'", testAwsLogsMessage))))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			defer container.Delete(ctx, containerd.WithSnapshotCleanup) //nolint:errcheck // testing only
			// Create a new task from the container and start it
			task, err := container.NewTask(ctx, creator)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			defer task.Delete(ctx) //nolint:errcheck // testing only

			err = task.Start(ctx)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			statusC, err := task.Wait(ctx)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			// Waiting for the task to finish
			status := <-statusC
			code, _, err := status.Result()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(code).Should(gomega.Equal(uint32(0)))

			// Validating in AWS logs
			cwOutput, err := cwClient.GetLogEvents(context.TODO(), &cloudwatchlogs.GetLogEventsInput{
				LogStreamName: aws.String(testAwsLogsStream),
				LogGroupName:  aws.String(testAwsLogsGroup),
				Limit:         aws.Int32(1),
			})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(*cwOutput.Events[0].Message).Should(gomega.Equal(testAwsLogsMessage))
		})
	})
}
