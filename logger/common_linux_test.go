// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build unit
// +build unit

package logger

import (
	"bytes"
	"context"
	"os"
	"sync/atomic"
	"testing"

	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/stretchr/testify/require"
)

// TestTracingLogRouting verifies the expected number of bytes are read from the source/containers' pipe files and the
// expected number of bytes are sent to the destination/the log driver.
func TestTracingLogRouting(t *testing.T) {
	// Create a tmp file that used to mock the io pipe where the logger reads log
	// messages from.
	tmpIOSource, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer os.Remove(tmpIOSource.Name()) //nolint:errcheck // testing only
	var (
		testStdout bytes.Buffer
		testStderr bytes.Buffer
	)
	// Create two pipes for stdout and stderr.
	inputForStdout := "1234567890\n"
	countOfNewLinesForStdout := 1
	_, err = testStdout.WriteString(inputForStdout)
	require.NoError(t, err)
	inputForStderr := "123 456 789 0\n123 456 789 0\n123 456 789 0\n"
	countOfNewLinesForStderr := 3
	_, err = testStderr.WriteString(inputForStderr)
	require.NoError(t, err)
	// Create a tmp file that used in customized dummy log function where the
	// logger sends log messages to.
	tmpDest, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer os.Remove(tmpDest.Name()) //nolint:errcheck // testing only
	logDestinationFileName = tmpDest.Name()

	l := &Logger{
		Info:              &dockerlogger.Info{},
		Stream:            &dummyClient{t},
		bufferSizeInBytes: DefaultBufSizeInBytes,
		maxReadBytes:      defaultMaxReadBytes,
		Stdout:            &testStdout,
		Stderr:            &testStderr,
	}
	err = l.Start(
		context.TODO(),
		&dummyCleanupTime,
		func() error { return nil },
	)
	require.NoError(t, err)

	require.Equal(t, uint64(len(inputForStdout)+len(inputForStderr)), atomic.LoadUint64(&bytesReadFromSrc))
	// Exclude the new line characters because they will be removed when sending logs to the log driver.
	require.Equal(t,
		uint64(len(inputForStdout)+len(inputForStderr)-countOfNewLinesForStdout-countOfNewLinesForStderr),
		atomic.LoadUint64(&bytesSentToDst))
	require.Equal(t,
		uint64(countOfNewLinesForStdout+countOfNewLinesForStderr), atomic.LoadUint64(&numberOfNewLineChars))
}
