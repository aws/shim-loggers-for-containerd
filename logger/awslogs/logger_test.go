// +build unit

package awslogs

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aws/shim-loggers-for-containerd/logger"
	mock_awslogs "github.com/aws/shim-loggers-for-containerd/logger/awslogs/mocks"

	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/golang/mock/gomock"
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
	testContainerID         = "test-container-id"
	testContainerName       = "test-container-name"
	testLogDriver           = "awslogs"
	testMode                = "blocking"
	maxRetries              = 3

	testErrMsg = "test error message"
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
	globalArgs = &logger.GlobalArgs{
		ContainerID:   testContainerID,
		ContainerName: testContainerName,
		LogDriver:     testLogDriver,
		Mode:          testMode,
	}

	dummyLogMsg            = []byte("test log message")
	dummyTime              = time.Date(2020, time.January, 14, 01, 59, 0, 0, time.UTC)
	logDestinationFileName string
)

// dummyClient is only used for testing. It owns the customized Log function used in
// TestSendLogs test case as we do not need the functionality that the actual Log function
// is doing inside the test. Mock Log function is not enough here as there does not exist a
// better way to verify what happened in the TestSendLogs test, which has a goroutine.
type dummyClient struct{}

// Log implements customized workflow used for testing purpose.
// This is only trigger in TestSendLogs test case. It writes current log message to the end of
// tmp test file, which makes sure the function itself accepts and "logging" the message
// correctly.
func (d *dummyClient) Log(msg *dockerlogger.Message) error {
	_, err := os.Stat(logDestinationFileName)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(logDestinationFileName, os.O_APPEND|os.O_RDWR, 0644)
	defer f.Close()
	if err != nil {
		return errors.Wrapf(err,
			"unable to open file %s to record log message", logDestinationFileName)
	}
	f.Write(msg.Line)

	return nil
}

// TestLogWithRetry tests function LogWithRetry does not retry on success or retries
// on error.
func TestLogWithRetry(t *testing.T) {
	t.Run("DoesNotRetry", testLogWithRetryDoesNotRetry)
	t.Run("WithError", testLogWithRetryWithError)
}

// testLogWithRetryDoesNotRetry tests LogWithRetry function did not retry on no error
// returned from Log function.
func testLogWithRetryDoesNotRetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStream := mock_awslogs.NewMockclient(ctrl)
	l := &logDriver{
		info:   &dockerlogger.Info{},
		stream: mockStream,
	}
	mockStream.EXPECT().Log(gomock.Any()).Return(nil).Times(1)
	err := l.LogWithRetry(dummyLogMsg, dummyTime)
	require.NoError(t, err)
}

// testLogWithRetryWithError tests LogWithRetry function retries on error returned from
// Log function.
func testLogWithRetryWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStream := mock_awslogs.NewMockclient(ctrl)
	l := &logDriver{
		info:   &dockerlogger.Info{},
		stream: mockStream,
	}
	mockStream.EXPECT().Log(gomock.Any()).Return(errors.New(testErrMsg)).Times(maxRetries)
	expectErrMsg := fmt.Sprintf("sending container logs to cloudwatch has been retried for %d times: %s",
		maxRetries, testErrMsg)
	err := l.LogWithRetry(dummyLogMsg, dummyTime)
	require.Error(t, err)
	require.Equal(t, expectErrMsg, err.Error())
}

// TestSendLogs tests sendLogs goroutine that gets log message from mock io pipe and sends
// to mock destination. In this test case, the source and destination are both tmp files that
// read from and write to inside the customized Log function.
func TestSendLogs(t *testing.T) {
	l := &logDriver{
		info:   &dockerlogger.Info{},
		stream: &dummyClient{},
	}
	// Create a tmp file that used to mock the io pipe where the logger reads log
	// messages from.
	tmpIOSource, err := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(tmpIOSource.Name())
	require.NoError(t, err)
	var expectedSize int64
	lines := []string{
		"First line to write",
		"Second line to write",
	}
	for _, line := range lines {
		expectedSize += int64(len(line))
		tmpIOSource.WriteString(line)
	}

	// Create a tmp file that used to inside customized dummy Log function where the
	// logger sends log messages to.
	tmpDest, err := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(tmpDest.Name())
	logDestinationFileName = tmpDest.Name()
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	f, err := os.Open(tmpIOSource.Name())
	defer f.Close()
	require.NoError(t, err)
	go l.sendLogs(f, &wg)
	wg.Wait()

	// Make sure the new scanned log message has been written to the tmp file by sendLogs
	// goroutine.
	logDestinationInfo, err := os.Stat(logDestinationFileName)
	require.NoError(t, err)
	require.Equal(t, expectedSize, logDestinationInfo.Size())
}

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

// TestNewInfo tests if newInfo function creates logger info correctly.
func TestNewInfo(t *testing.T) {
	config := getAWSLogsConfig(args)
	info := newInfo(testContainerID, testContainerName, WithConfig(config))
	require.Equal(t, config, info.Config)
}
