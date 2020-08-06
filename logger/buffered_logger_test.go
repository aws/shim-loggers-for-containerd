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
	"fmt"
	"testing"

	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/stretchr/testify/require"
)

const (
	testBufferSize = 100
)

var (
	messages = []*dockerlogger.Message{
		{Line: []byte("line1"), Timestamp: dummyTime},
		{Line: []byte("line2"), Timestamp: dummyTime},
		{Line: []byte("line3"), Timestamp: dummyTime},
		{Line: []byte("testLine4"), Timestamp: dummyTime},
	}
)

// testEnqueue tests Enqueue operation without error and gets used
// as initialization of buffer in Dequeue and Flush tests.
func testEnqueue(t *testing.T) *ringBuffer {
	lb := newLoggerBuffer(testBufferSize)
	require.Equal(t, testBufferSize, lb.maxSizeInBytes)
	require.Equal(t, 0, lb.curSizeInBytes)

	var (
		err                   error
		expectedCurBufferSize int
	)
	for _, msg := range messages {
		expectedCurBufferSize += len(msg.Line)
		err = lb.Enqueue(msg)
		require.NoError(t, err)
	}
	require.Len(t, lb.queue, len(messages))
	require.Equal(t, expectedCurBufferSize, lb.curSizeInBytes)

	return lb
}

// TestLogBufferEnqueueDequeue tests dequeue operations from
// buffer
func TestLogBufferEnqueueDequeue(t *testing.T) {
	lb := testEnqueue(t)
	queueLen := len(lb.queue)
	for i := 0; i < queueLen; i++ {
		fmt.Println(i)
		msg, err := lb.Dequeue()
		require.NoError(t, err)
		require.Equal(t, messages[i], msg)
	}
	require.Len(t, lb.queue, 0)
}

// TestLogBufferEnqueueFlush tests flush messages from buffer
func TestLogBufferEnqueueFlush(t *testing.T) {
	lb := testEnqueue(t)
	flushedMsg := lb.Flush()
	require.Len(t, lb.queue, 0)
	require.Equal(t, messages, flushedMsg)
}
