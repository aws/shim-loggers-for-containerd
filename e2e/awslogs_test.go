// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build e2e

package e2e

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	smithy "github.com/aws/smithy-go"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/containerd/containerd/cio"
	"github.com/google/uuid"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	awslogsCredentialsEndpointKey = "--awslogs-credentials-endpoint" //nolint:gosec // not credentials
	awslogsRegionKey              = "--awslogs-region"
	awslogsStreamKey              = "--awslogs-stream"
	awslogsGroupKey               = "--awslogs-group"
	awslogsEndpointKey            = "--awslogs-endpoint"
	awslogsCreateGroupKey         = "--awslogs-create-group"
	awslogsCreateStreamKey        = "--awslogs-create-stream"
	awslogsMultilinePatternKey    = "--awslogs-multiline-pattern"
	awslogsDatetimeFormatKey      = "--awslogs-datetime-format"
	awslogsLogFormatKey           = "--awslogs-format"
	testEcsLocalEndpointPort      = "51679"
	testAwslogsCredentialEndpoint = ":" + testEcsLocalEndpointPort + "/creds"
	testAwslogsRegion             = "us-east-1"
	testAwslogsStream             = "test-stream"
	nonExistentAwslogsStream      = "non-existent-stream"
	testAwslogsGroup              = "test-shim-logger"
	nonExistentAwslogsGroup       = "non-existent-group"
	testAwslogsEndpoint           = "http://localhost.localstack.cloud:4566" // Recommended endpoint:
	testAwslogsJSONEmf            = "json/emf"
	//nolint:lll // url
	// https://docs.localstack.cloud/getting-started/faq/#is-using-localhostlocalstackcloud4566-to-set-as-the-endpoint-for-aws-services-recommended
	testAwslogsMultilinePattern = "TEST"
	testAwslogsDatetimeFormat   = "\\[%b %d, %Y %H:%M:%S\\]"
)

type cloudWatchEndpointResolver struct{}

func (cloudWatchEndpointResolver) ResolveEndpoint(
	_ context.Context,
	_ cloudwatchlogs.EndpointParameters,
) (smithyendpoints.Endpoint, error) {
	uri, _ := url.Parse(testAwslogsEndpoint)
	return smithyendpoints.Endpoint{
		URI: *uri,
	}, nil
}

var testAwslogs = func() {
	// These tests are run in serial because we only define one log driver instance.
	ginkgo.Describe("awslogs shim logger", ginkgo.Serial, func() {
		var cwClient *cloudwatchlogs.Client
		ginkgo.BeforeEach(func() {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(testAwslogsRegion))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			cwClient = cloudwatchlogs.NewFromConfig(cfg, func(opts *cloudwatchlogs.Options) {
				opts.EndpointResolverV2 = cloudWatchEndpointResolver{}
			})
			deleteLogGroup(cwClient, testAwslogsGroup)
			deleteLogGroup(cwClient, nonExistentAwslogsGroup)
			_, err = cwClient.CreateLogGroup(context.TODO(), &cloudwatchlogs.CreateLogGroupInput{
				LogGroupName: aws.String(testAwslogsGroup),
			})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			_, err = cwClient.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
				LogGroupName:  aws.String(testAwslogsGroup),
				LogStreamName: aws.String(testAwslogsStream),
			})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.AfterEach(func() {
			deleteLogGroup(cwClient, testAwslogsGroup)
			deleteLogGroup(cwClient, nonExistentAwslogsGroup)
		})
		ginkgo.It("should send logs to awslogs log driver with existing log group and non-existent log stream "+
			"when the configs are default", func() {
			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwslogsGroup,
				awslogsStreamKey:              nonExistentAwslogsStream,
				awslogsEndpointKey:            testAwslogsEndpoint,
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			err = validateTestLogsInAwslogs(cwClient, testAwslogsGroup, nonExistentAwslogsStream, []string{testLog})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.It("should send logs to awslogs log driver with respecting datatime format and ignoring multiline pattern "+
			"when both datetime format and multiline pattern are set", func() {
			firstLineLog := "[May 01, 2017 19:00:01] " + testLogPrefix + uuid.New().String()
			secondLineLog := "[May 01, 2017 19:00:04] " + testLogPrefix + uuid.New().String()
			thirdLineLog := fmt.Sprintf("%s %s", testAwslogsMultilinePattern, testLogPrefix+uuid.New().String())
			fourthLineLog := fmt.Sprintf("%s %s", testAwslogsMultilinePattern, testLogPrefix+uuid.New().String())
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwslogsGroup,
				awslogsStreamKey:              nonExistentAwslogsStream,
				awslogsEndpointKey:            testAwslogsEndpoint,
				awslogsMultilinePatternKey:    "^" + testAwslogsMultilinePattern,
				awslogsDatetimeFormatKey:      testAwslogsDatetimeFormat,
			}
			creator := cio.BinaryIO(*Binary, args)
			// The last matched line cannot be logged with multiline pattern. Append a pattern for now.
			// TODO: Investigate and fix. https://github.com/aws/shim-loggers-for-containerd/issues/78
			err := SendTestLogByContainerd(creator, fmt.Sprintf("%s\n%s\n%s\n%s\n%s", firstLineLog,
				secondLineLog, thirdLineLog, fourthLineLog, "[May 01, 2017 19:00:05]"))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			err = validateTestLogsInAwslogs(cwClient, testAwslogsGroup, nonExistentAwslogsStream, []string{fmt.Sprintf("%s\n", firstLineLog),
				fmt.Sprintf("%s\n%s\n%s\n", secondLineLog, thirdLineLog, fourthLineLog)})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.It("should send logs to awslogs log driver with respecting multiline pattern "+
			"when multiline pattern is set", func() {
			firstLineLog := fmt.Sprintf("%s %s", testAwslogsMultilinePattern, testLogPrefix+uuid.New().String())
			secondLineLog := fmt.Sprintf("%s %s", testAwslogsMultilinePattern, testLogPrefix+uuid.New().String())
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwslogsGroup,
				awslogsStreamKey:              nonExistentAwslogsStream,
				awslogsEndpointKey:            testAwslogsEndpoint,
				awslogsMultilinePatternKey:    "^" + testAwslogsMultilinePattern,
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, fmt.Sprintf("%s\n%s\n%s", firstLineLog, secondLineLog, testAwslogsMultilinePattern))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			err = validateTestLogsInAwslogs(cwClient, testAwslogsGroup, nonExistentAwslogsStream,
				[]string{fmt.Sprintf("%s\n", firstLineLog), fmt.Sprintf("%s\n", secondLineLog)})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.It("should send logs to awslogs log driver with respecting datatime format "+
			"when datetime format is set", func() {
			testLog := testLogPrefix + uuid.New().String()
			firstLineLog := "[May 01, 2017 19:00:01] " + testLog
			secondLineLog := "[May 01, 2017 19:00:04] " + testLog
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwslogsGroup,
				awslogsStreamKey:              nonExistentAwslogsStream,
				awslogsEndpointKey:            testAwslogsEndpoint,
				awslogsDatetimeFormatKey:      testAwslogsDatetimeFormat,
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, fmt.Sprintf("%s\n%s\n%s", firstLineLog, secondLineLog, "[May 01, 2017 19:00:05]"))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			err = validateTestLogsInAwslogs(cwClient, testAwslogsGroup, nonExistentAwslogsStream,
				[]string{fmt.Sprintf("%s\n", firstLineLog), fmt.Sprintf("%s\n", secondLineLog)})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.It("should send logs to awslogs log driver with existing log group and non-existent log stream "+
			"when createGroup is false and createStream is true", func() {
			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwslogsGroup,
				awslogsStreamKey:              nonExistentAwslogsStream,
				awslogsEndpointKey:            testAwslogsEndpoint,
				awslogsCreateGroupKey:         "false",
				awslogsCreateStreamKey:        "true",
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			err = validateTestLogsInAwslogs(cwClient, testAwslogsGroup, nonExistentAwslogsStream, []string{testLog})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.It("should fail to send logs to awslogs log driver with non-existent log group "+
			"when createGroup is false", func() {
			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               nonExistentAwslogsGroup,
				awslogsStreamKey:              nonExistentAwslogsGroup,
				awslogsEndpointKey:            testAwslogsEndpoint,
				awslogsCreateGroupKey:         "false",
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.Equal(containerdTaskExitNonZeroMessage))
		})
		ginkgo.It("should silently fail to send logs to awslogs log driver with existing log group and non-existent log stream "+
			"when createGroup is false and createStream is false", func() {
			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwslogsGroup,
				awslogsStreamKey:              nonExistentAwslogsStream,
				awslogsEndpointKey:            testAwslogsEndpoint,
				awslogsCreateGroupKey:         "false",
				awslogsCreateStreamKey:        "false",
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			err = validateTestLogsInAwslogs(cwClient, testAwslogsGroup, nonExistentAwslogsStream, []string{testLog})
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})
		ginkgo.It("should send logs to awslogs log driver with non-existent log group and non-existent log stream "+
			"when createGroup is true and createStream is true", func() {
			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               nonExistentAwslogsGroup,
				awslogsStreamKey:              nonExistentAwslogsStream,
				awslogsEndpointKey:            testAwslogsEndpoint,
				awslogsCreateGroupKey:         "true",
				awslogsCreateStreamKey:        "true",
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			err = validateTestLogsInAwslogs(cwClient, nonExistentAwslogsGroup, nonExistentAwslogsStream, []string{testLog})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.It("should send logs to awslogs log driver with existing log group and log stream when createGroup "+
			"is false and createStream is false", func() {
			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwslogsGroup,
				awslogsStreamKey:              testAwslogsStream,
				awslogsEndpointKey:            testAwslogsEndpoint,
				awslogsCreateGroupKey:         "false",
				awslogsCreateStreamKey:        "false",
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			err = validateTestLogsInAwslogs(cwClient, testAwslogsGroup, testAwslogsStream, []string{testLog})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		ginkgo.It("should set the header 'x-amzn-logs-format' to have value 'json/emf' "+
			"when LogsFormatHeader is configured with 'json/emf'", func() {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gomega.Expect(r.Header.Get("x-amzn-logs-format")).To(gomega.Equal("json/emf"))
				w.WriteHeader(http.StatusOK)
			}))
			defer mockServer.Close()

			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:              AwslogsDriverName,
				ContainerIDKey:                TestContainerID,
				ContainerNameKey:              TestContainerName,
				awslogsCredentialsEndpointKey: testAwslogsCredentialEndpoint,
				awslogsRegionKey:              testAwslogsRegion,
				awslogsGroupKey:               testAwslogsGroup,
				awslogsStreamKey:              nonExistentAwslogsStream,
				awslogsLogFormatKey:           testAwslogsJSONEmf,
				awslogsEndpointKey:            mockServer.URL,
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
	})
}

func validateTestLogsInAwslogs(client *cloudwatchlogs.Client, logGroupName string, logStreamName string, testLogs []string) error {
	cwOutput, err := client.GetLogEvents(context.TODO(), &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	})
	if err != nil {
		return err
	}
	co := *cwOutput
	if len(co.Events) != len(testLogs) {
		return fmt.Errorf("the number of test log lines are not matching")
	}
	for i := 0; i < len(testLogs); i++ {
		if *co.Events[i].Message != testLogs[i] {
			return fmt.Errorf("test log messages are not matching at %d line", i)
		}
	}
	return nil
}

func deleteLogGroup(client *cloudwatchlogs.Client, logGroupName string) {
	_, err := client.DeleteLogGroup(context.TODO(), &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(logGroupName),
	})
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			gomega.Expect(oe.Unwrap().Error()).Should(
				gomega.ContainSubstring("ResourceNotFoundException: The specified log group does not exist"))
		}
	} else {
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	}
}
