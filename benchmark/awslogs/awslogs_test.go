// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build benchmark

package awslogs

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/shim-loggers-for-containerd/e2e"
	"github.com/containerd/containerd/cio"
)

const (
	awslogsCredentialsEndpointKey = "--awslogs-credentials-endpoint" //nolint:gosec // not credentials
	awslogsRegionKey              = "--awslogs-region"
	awslogsStreamKey              = "--awslogs-stream"
	awslogsGroupKey               = "--awslogs-group"
	awslogsEndpointKey            = "--awslogs-endpoint"
	testEcsLocalEndpointPort      = "51679"
	testAwslogsCredentialEndpoint = ":" + testEcsLocalEndpointPort + "/creds"
	testAwslogsRegion             = "us-east-1"
	testAwslogsStream             = "test-stream"
	nonExistentAwslogsStream      = "non-existent-stream"
	testAwslogsGroup              = "test-shim-logger"
	testAwslogsEndpoint           = "http://localhost.localstack.cloud:4566" // Recommended endpoint:
)

var (
	// Binary is the path the binary of the shim loggers for containerd.
	Binary = flag.String("binary", "", "the binary of shim loggers for containerd")
)

func BenchmarkAwslogs(b *testing.B) {
	var cwClient *cloudwatchlogs.Client
	// Reference to set up Go client for aws local stack: https://docs.localstack.cloud/user-guide/integrations/sdks/go/.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           testAwslogsEndpoint,
			SigningRegion: testAwslogsRegion,
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(testAwslogsRegion),
		config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		b.Fatal(err)
	}
	cwClient = cloudwatchlogs.NewFromConfig(cfg)
	_, err = cwClient.CreateLogGroup(context.TODO(), &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(testAwslogsGroup),
	})
	if err != nil {
		b.Fatal(err)
	}
	_, err = cwClient.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(testAwslogsGroup),
		LogStreamName: aws.String(testAwslogsStream),
	})
	if err != nil {
		b.Fatal(err)
	}
	testLog := strings.Repeat("a", 1024)
	args := map[string]string{
		e2e.LogDriverTypeKey:          e2e.AwslogsDriverName,
		e2e.ContainerIDKey:            e2e.TestContainerID,
		e2e.ContainerNameKey:          e2e.TestContainerName,
		awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
		awslogsRegionKey:              testAwslogsRegion,
		awslogsGroupKey:               testAwslogsGroup,
		awslogsStreamKey:              nonExistentAwslogsStream,
		awslogsEndpointKey:            testAwslogsEndpoint,
	}
	creator := cio.BinaryIO(*Binary, args)
	err = e2e.SendTestLogByContainerd(creator, testLog)
	if err != nil {
		b.Fatal(err)
	}
}
