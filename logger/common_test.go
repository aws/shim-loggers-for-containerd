// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build unit
// +build unit

package logger

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

const (
	maxRetries        = 3
	testErrMsg        = "test error message"
	testContainerID   = "test-container-id"
	testContainerName = "test-container-name"
)

var (
	dummySource            = "stdout"
	dummyTime              = time.Date(2020, time.January, 14, 01, 59, 0, 0, time.UTC)
	dummyCleanupTime       = 2 * time.Second
	logDestinationFileName string
)

// dummyClient is only used for testing. It owns the customized Log function used in
// TestSendLogs test case as we do not need the functionality that the actual Log function
// is doing inside the test. Mock Log function is not enough here as there does not exist a
// better way to verify what happened in the TestSendLogs test, which has a goroutine.
type dummyClient struct {
	t *testing.T

	// logCallFailAtIndex is used to signal when the `Log` method will fail. Use 0 to disable it.
	logCallFailAtIndex uint64
	// logCallSucceedAtIndex is used to signal when the `Log` method will succeed.
	logCallSucceedAtIndex uint64
	logCallTimes          uint64
}

// Log implements customized workflow used for testing purpose.
// This is only trigger in TestSendLogs test case. It writes current log message to the end of
// tmp test file, which makes sure the function itself accepts and "logging" the message
// correctly.
func (d *dummyClient) Log(msg *dockerlogger.Message) error {
	logCallTimes := atomic.LoadUint64(&(d.logCallTimes))
	logCallFailsAt := atomic.LoadUint64(&d.logCallFailAtIndex)
	logCallSucceedAt := atomic.LoadUint64(&d.logCallSucceedAtIndex)

	if logCallTimes >= logCallFailsAt && logCallTimes < logCallSucceedAt && logCallFailsAt != 0 {
		currentLogCallTimes := atomic.AddUint64(&d.logCallTimes, 1)
		return fmt.Errorf("fail `Log` intentionally at logCallTimes %d", currentLogCallTimes)
	}

	atomic.AddUint64(&d.logCallTimes, 1)

	var b []byte
	_, err := os.Stat(logDestinationFileName)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(logDestinationFileName, os.O_APPEND|os.O_RDWR, 0644) //nolint:gosec // testing only
	if err != nil {
		return fmt.Errorf("unable to open file %s to record log message: %w", logDestinationFileName, err)
	}
	defer f.Close() //nolint:errcheck // testing only
	b, err = json.Marshal(msg)
	require.NoError(d.t, err)
	f.Write(b)            //nolint:errcheck,gosec // testing only
	f.Write([]byte{'\n'}) //nolint:errcheck,gosec // testing only

	return nil
}

func checkLogFile(t *testing.T, fileName string, expectedNumLines int,
	expectedPartialOrdinalSequence []int) {
	var (
		msg                dockerlogger.Message
		line               string
		lastPartialID      string
		lastPartialOrdinal int
	)
	file, err := os.Open(fileName) //nolint:gosec // testing only
	require.NoError(t, err)
	defer file.Close() //nolint:errcheck // testing only

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		line = scanner.Text()
		err = json.Unmarshal([]byte(line), &msg)
		require.NoError(t, err)
		if len(expectedPartialOrdinalSequence) > 0 && lines < len(expectedPartialOrdinalSequence) {
			// check partial fields
			// if the message is partial it will have a PLogMetaData
			require.NotNil(t, msg.PLogMetaData)
			require.Equal(t, expectedPartialOrdinalSequence[lines], msg.PLogMetaData.Ordinal)
			if msg.PLogMetaData.Ordinal < lastPartialOrdinal {
				// new split message so new partial ID
				require.NotEqual(t, lastPartialID, msg.PLogMetaData.ID)
			} else if msg.PLogMetaData.Ordinal > 1 {
				// this partial ID should be same as last ID
				require.Equal(t, lastPartialID, msg.PLogMetaData.ID)
			}
			lastPartialID = msg.PLogMetaData.ID
		}
		lines++
	}
	require.Equal(t, expectedNumLines, lines)

	err = scanner.Err()
	require.NoError(t, err)
}

// TestSendLogs tests sendLogs goroutine that gets log message from mock io pipe and sends
// to mock destination. In this test case, the source and destination are both tmp files that
// read from and write to inside the customized Log function.
func TestSendLogs(t *testing.T) {
	for _, tc := range []struct {
		testName                       string
		bufferSizeInBytes              int
		maxReadBytes                   int
		logMessages                    []string
		expectedNumOfLines             int
		expectedPartialOrdinalSequence []int
	}{
		{
			testName:          "general case",
			bufferSizeInBytes: 100,
			maxReadBytes:      80, // Larger than the sum of sizes of two log messages.
			logMessages: []string{
				"First line to write",
				"Second line to write",
			},
			expectedNumOfLines:             2,       // 2 messages stay as 2 messages
			expectedPartialOrdinalSequence: []int{}, // neither will be partial
		},
		{
			testName:          "long log message",
			bufferSizeInBytes: 8,
			maxReadBytes:      4,
			logMessages: []string{
				"First line to write", // Larger than buffer size.
			},
			expectedNumOfLines:             3, // One line 19 chars with 8 char buffer becomes 3 split messages
			expectedPartialOrdinalSequence: []int{1, 2, 3},
		},
		{
			testName:          "two long log messages",
			bufferSizeInBytes: 8,
			maxReadBytes:      4,
			logMessages: []string{
				"First line to write",  // 19 chars => 3 messages
				"Second line to write", // 20 chars => 3 messages
			},
			expectedNumOfLines:             6, // 3 + 3 = 6 total
			expectedPartialOrdinalSequence: []int{1, 2, 3, 1, 2, 3},
		},
	} {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			l := &Logger{
				Info: &dockerlogger.Info{},
				Stream: &dummyClient{
					t: t,
				},
				bufferSizeInBytes: tc.bufferSizeInBytes,
				maxReadBytes:      tc.maxReadBytes,
			}
			// Create a tmp file that used to mock the io pipe where the logger reads log
			// messages from.
			tmpIOSource, err := os.CreateTemp("", "")
			require.NoError(t, err)
			defer os.Remove(tmpIOSource.Name()) //nolint:errcheck // testing only
			var (
				testPipe bytes.Buffer
			)
			for _, logMessage := range tc.logMessages {
				_, err := testPipe.WriteString(logMessage + "\n")
				require.NoError(t, err)
			}

			// Create a tmp file that used to inside customized dummy Log function where the
			// logger sends log messages to.
			tmpDest, err := os.CreateTemp("", "")
			require.NoError(t, err)
			defer os.Remove(tmpDest.Name()) //nolint:errcheck // testing only
			logDestinationFileName = tmpDest.Name()

			var errGroup errgroup.Group
			errGroup.Go(func() error {
				return l.sendLogs(context.TODO(), &testPipe, dummySource, &dummyCleanupTime)
			})
			err = errGroup.Wait()
			require.NoError(t, err)

			// Make sure the new scanned log message has been written to the tmp file by sendLogs
			// goroutine.
			logDestinationInfo, err := os.Stat(logDestinationFileName)
			require.NoError(t, err)
			require.NotZero(t, logDestinationInfo.Size())

			checkLogFile(t, logDestinationFileName, tc.expectedNumOfLines, tc.expectedPartialOrdinalSequence)
		})
	}
}

// TestNewInfo tests if NewInfo function creates logger info correctly.
func TestNewInfo(t *testing.T) {
	config := map[string]string{
		"configKey1": "configVal1",
		"configKey2": "configVal2",
		"configKey3": "configVal3",
	}
	info := NewInfo(testContainerID, testContainerName, WithConfig(config))
	require.Equal(t, config, info.Config)
}

// TestPipeNotBroken verified the pipe will NOT be broken even if sometimes the call
// to the log driver fails.
func TestPipeNotBroken(t *testing.T) {
	logCallFailAtIndex := uint64(1)
	logCallSucceedAtIndex := uint64(3)
	logCallTimes := uint64(0)

	l := &Logger{
		Info: &dockerlogger.Info{},
		Stream: &dummyClient{
			t:                     t,
			logCallFailAtIndex:    logCallFailAtIndex,
			logCallSucceedAtIndex: logCallSucceedAtIndex,
			logCallTimes:          logCallTimes,
		},
		bufferSizeInBytes: 8,
		maxReadBytes:      4,
	}

	// 60 chars with 8 char buffer becomes 8 split messages.
	msgFromIOSource := "First line to write Second line to write Third line to write"

	// Create a tmp file that used to mock the io pipe where the logger reads log
	// messages from.
	tmpIOSource, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer os.Remove(tmpIOSource.Name()) //nolint:errcheck // testing only
	var testPipe bytes.Buffer
	_, err = testPipe.WriteString(msgFromIOSource + "\n")
	require.NoError(t, err)

	// Create a tmp file that used to inside customized dummy Log function where the
	// logger sends log messages to.
	tmpDest, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer os.Remove(tmpDest.Name()) //nolint:errcheck // testing only
	logDestinationFileName = tmpDest.Name()

	var errGroup errgroup.Group
	errGroup.Go(func() error {
		return l.sendLogs(context.TODO(), &testPipe, dummySource, &dummyCleanupTime)
	})
	err = errGroup.Wait()
	require.NoError(t, err)

	// Verify that the log destination received partial msg only.
	file, err := os.Open(logDestinationFileName) //nolint:gosec // testing only
	require.NoError(t, err)
	defer file.Close() //nolint:errcheck // testing only

	scanner := bufio.NewScanner(file)
	lines := 0
	var receivedMsgInDst strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		var msg dockerlogger.Message
		err = json.Unmarshal([]byte(line), &msg)
		t.Logf("Received msg: %v", string(msg.Line))
		receivedMsgInDst.WriteString(string(msg.Line))
		require.NoError(t, err)
		lines++
	}

	// 6 lines will be received by log driver as we fail 2 calls intentionally.
	require.Equal(t, 6, lines)
	require.Equal(t, receivedMsgInDst.String(), "First lind line to write Third line to write")

	err = scanner.Err()
	require.NoError(t, err)
}
