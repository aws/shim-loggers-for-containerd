// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// +build unit

package logger

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/pkg/errors"
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
	dummyLogMsg            = []byte("test log message")
	dummySource            = "stdout"
	dummyTime              = time.Date(2020, time.January, 14, 01, 59, 0, 0, time.UTC)
	dummyCleanupTime       = time.Duration(2 * time.Second)
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
	if err != nil {
		return errors.Wrapf(err,
			"unable to open file %s to record log message", logDestinationFileName)
	}
	defer f.Close()
	f.Write(msg.Line)

	return nil
}

// TestSendLogs tests sendLogs goroutine that gets log message from mock io pipe and sends
// to mock destination. In this test case, the source and destination are both tmp files that
// read from and write to inside the customized Log function.
func TestSendLogs(t *testing.T) {

	for _, tc := range []struct {
		testName           string
		bufferSizeInBytes  int
		maxReadBytes       int
		logMessages        []string
		expectedNumOfLines int
	}{
		{
			testName:          "general case",
			bufferSizeInBytes: 100,
			maxReadBytes:      80, // Larger than the sum of sizes of two log messages.
			logMessages: []string{
				"First line to write",
				"Second line to write",
			},
			expectedNumOfLines: 2,
		},
		{
			testName:          "long log message",
			bufferSizeInBytes: 8,
			maxReadBytes:      4,
			logMessages: []string{
				"First line to write", // Larger than buffer size.
			},
			expectedNumOfLines: 2, // Should be split to 2 lines in destination.
		},
		{
			testName:          "partial log message",
			bufferSizeInBytes: 100,
			maxReadBytes:      4,
			logMessages: []string{
				"First line to write", // Larger than maximum allowed read bytes from pipe.
				"Second line to write",
			},
			expectedNumOfLines: 2, // Should still be 2 lines in destination.
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			l := &Logger{
				Info:              &dockerlogger.Info{},
				Stream:            &dummyClient{},
				bufferSizeInBytes: tc.bufferSizeInBytes,
				maxReadBytes:      tc.maxReadBytes,
			}
			// Create a tmp file that used to mock the io pipe where the logger reads log
			// messages from.
			tmpIOSource, err := ioutil.TempFile("", "")
			require.NoError(t, err)
			defer os.Remove(tmpIOSource.Name())
			var (
				expectedSize int64
				testPipe     bytes.Buffer
			)
			for _, logMessage := range tc.logMessages {
				expectedSize += int64(len([]rune(logMessage)))
				_, err := testPipe.WriteString(logMessage)
				require.NoError(t, err)
			}

			// Create a tmp file that used to inside customized dummy Log function where the
			// logger sends log messages to.
			tmpDest, err := ioutil.TempFile(os.TempDir(), "")
			require.NoError(t, err)
			defer os.Remove(tmpDest.Name())
			logDestinationFileName = tmpDest.Name()

			var errGroup errgroup.Group
			errGroup.Go(func() error {
				return l.sendLogs(context.TODO(), &testPipe, dummySource, -1, -1, &dummyCleanupTime)
			})
			err = errGroup.Wait()
			require.NoError(t, err)

			// Make sure the new scanned log message has been written to the tmp file by sendLogs
			// goroutine.
			logDestinationInfo, err := os.Stat(logDestinationFileName)
			require.NoError(t, err)
			require.Equal(t, expectedSize, logDestinationInfo.Size())
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
