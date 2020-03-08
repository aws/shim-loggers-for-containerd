// +build unit

package awslogs

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	testGroup               = "test-group"
	testRegion              = "test-region"
	testStream              = "test-stream"
	testCredentialsEndpoint = "test-credential-endpoints"
	testCreateGroup         = "true"
	testMultilinePattern    = "test-multiline-pattern"
	testDatetimeFormat      = "test-datetime-format"
)

var (
	args = &Args{
		Group:               testGroup,
		Region:              testRegion,
		Stream:              testStream,
		CredentialsEndpoint: testCredentialsEndpoint,
		CreateGroup:         testCreateGroup,
		MultilinePattern:    testMultilinePattern,
		DatetimeFormat:      testDatetimeFormat,
	}
)

// TestGetAWSLogsConfig tests if getAWSLogsConfig function converts log config
// maps correctly.
func TestGetAWSLogsConfig(t *testing.T) {
	expectedConfig := map[string]string{
		GroupKey:               testGroup,
		RegionKey:              testRegion,
		StreamKey:              testStream,
		CredentialsEndpointKey: testCredentialsEndpoint,
		CreateGroupKey:         testCreateGroup,
		MultilinePatternKey:    testMultilinePattern,
		DatetimeFormatKey:      testDatetimeFormat,
	}

	config := getAWSLogsConfig(args)
	require.Equal(t, expectedConfig, config)
}

// TestCreateStreamWithRetry tests if function CreateStreamWithRetry does try on retiable error
// or not retry on no errors or non-retriable errors.
func TestCreateStreamWithRetry(t *testing.T) {
	t.Run("CreateStreamRetriesOnRetriableError", testCreateStreamRetriesOnRetriableErr)
	t.Run("CreateStreamDoesNotRetryOnNoError", testCreateStreamDoesNotRetryOnNoErr)
	t.Run("CreateStreamDoesNotRetryOnNonRetriableError", testCreateStreamDoesNotRetryOnNonRetriableErr)
}

// mockDockerNewLogger is customized mock function of docker New API for creating Logger. It returns
// error with messages for different test cases.
func mockDockerNewLogger(info dockerlogger.Info) (dockerlogger.Logger, error) {
	switch info.ContainerName {
	case "NoErr":
		return nil, nil
	case "RetriableErr":
		opAbortedError := awserr.New(cloudwatchlogs.ErrCodeOperationAbortedException,
			"operation aborted", errors.New("operation aborted"))
		return nil, errors.Wrap(opAbortedError, "retriable error wrapper")
	case "NoRetriableErr":
		return nil, errors.New("not a retriable error")
	}

	return nil, nil
}

// testCreateStreamRetriesOnRetriableErr tests if createStreamWithRetry retries on retriable error.
func testCreateStreamRetriesOnRetriableErr(t *testing.T) {
	info := &dockerlogger.Info{
		ContainerName: "RetriableErr",
	}
	_, err := createStreamWithRetry(mockDockerNewLogger, info)
	require.Error(t, err)
	require.Contains(t, err.Error(), "creating log stream has been retried 5 time(s)")
}

// testCreateStreamDoesNotRetryOnNoErr tests if createStreamWithRetry does not try when no errors
// returned.
func testCreateStreamDoesNotRetryOnNoErr(t *testing.T) {
	info := &dockerlogger.Info{
		ContainerName: "NoErr",
	}
	_, err := createStreamWithRetry(mockDockerNewLogger, info)
	require.NoError(t, err)
}

// testCreateStreamDoesNotRetryOnNonRetriableErr tests if createStreamWithRetry does not try on
// non-retriable error when there's error returned when creating stream.
func testCreateStreamDoesNotRetryOnNonRetriableErr(t *testing.T) {
	info := &dockerlogger.Info{
		ContainerName: "NoRetriableErr",
	}
	_, err := createStreamWithRetry(mockDockerNewLogger, info)
	require.Error(t, err)
	require.NotContains(t, err.Error(), "creating log stream has been retried 5 time(s)")
}
