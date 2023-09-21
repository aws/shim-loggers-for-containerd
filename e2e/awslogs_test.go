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
	awslogEndpointKey             = "--awslogs-endpoint"
	testEcsLocalEndpointPort      = "51679"
	testAwslogsCredentialEndpoint = ":" + testEcsLocalEndpointPort + "/creds"
	testAwslogsRegion             = "us-east-1"
	testAwslogsStream             = "test-stream"
	testAwslogsGroup              = "test-shim-logger"
	testAwslogsEndpoint           = "localhost.localstack.cloud" // Recommended endpoint: https://docs.localstack.cloud/getting-started/faq/#is-using-localhostlocalstackcloud4566-to-set-as-the-endpoint-for-aws-services-recommended
)

var testAwslogs = func() {
	// These tests are run in serial because we only define one log driver instance.
	ginkgo.Describe("awslogs shim logger", ginkgo.Serial, func() {
		var cwClient *cloudwatchlogs.Client
		ginkgo.BeforeEach(func() {
			// Reference to set up Go client for aws local stack: https://docs.localstack.cloud/user-guide/integrations/sdks/go/.
			customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           testAwslogsEndpoint,
					SigningRegion: testAwslogsRegion,
				}, nil
			})
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(testAwslogsRegion), config.WithEndpointResolverWithOptions(customResolver))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			cwClient = cloudwatchlogs.NewFromConfig(cfg)
			cleanupAwslogs(cwClient, testAwslogsGroup, testAwslogsStream)
		})
		ginkgo.AfterEach(func() {
			cleanupAwslogs(cwClient, testAwslogsGroup, testAwslogsStream)
		})
		ginkgo.It("should send logs to awslogs log driver", func() {
			args := map[string]string{
				logDriverTypeKey:              awslogsDriverName,
				containerIdKey:                testContainerId,
				containerNameKey:              testContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwslogsGroup,
				awslogsStreamKey:              testAwslogsStream,
				awslogEndpointKey:             testAwslogsEndpoint,
			}
			creator := cio.BinaryIO(*Binary, args)
			sendTestLogByContainerd(creator, testLog)
			validateTestLogsInAwslogs(cwClient, testAwslogsGroup, testAwslogsStream, testLog)
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
