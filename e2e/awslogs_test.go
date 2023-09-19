// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/containerd/containerd/cio"
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
)

var testAwslogs = func() {
	// These tests are run in serial because we only define one log driver instance.
	ginkgo.Describe("awslogs shim logger", ginkgo.Serial, func() {
		var cwClient *cloudwatchlogs.Client
		ginkgo.BeforeEach(func() {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(testAwslogsRegion))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			cwClient = cloudwatchlogs.NewFromConfig(cfg)
			cleanupAwslogs(cwClient, testAwsLogsGroup, testAwsLogsStream)
		})
		ginkgo.AfterEach(func() {
			cleanupAwslogs(cwClient, testAwsLogsGroup, testAwsLogsStream)
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
			creator := cio.BinaryIO(*Binary, args)
			sendTestLogByContainerd(creator, testLog)
			validateTestLogsInAwslogs(cwClient, testAwsLogsGroup, testAwsLogsStream, testLog)
		})
	})
}

func validateTestLogsInAwslogs(client *cloudwatchlogs.Client, logGroupName string, logStreamName string, testLog string) {
	cwOutput, err := client.GetLogEvents(context.TODO(), &cloudwatchlogs.GetLogEventsInput{
		LogStreamName: aws.String(logGroupName),
		LogGroupName:  aws.String(logStreamName),
		Limit:         aws.Int32(1),
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(*cwOutput.Events[0].Message).Should(gomega.Equal(testLog))
}

func cleanupAwslogs(client *cloudwatchlogs.Client, logGroupName string, logStreamName string) {
	_, err := client.DeleteLogStream(context.TODO(), &cloudwatchlogs.DeleteLogStreamInput{
		LogStreamName: aws.String(logGroupName),
		LogGroupName:  aws.String(logStreamName),
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}
