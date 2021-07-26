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

package logger

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"

	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	expectedNumOfPipes = 2
	// This value is adopted from Docker:
	// https://github.com/moby/moby/blob/master/daemon/logger/ring.go#L140
	ringCap = 1000
)

// bufferedLogger is a wrapper of underlying log driver and an intermediate ring
// buffer between container pipes and underlying log driver.
type bufferedLogger struct {
	l           LogDriver
	buffer      *ringBuffer
	containerID string
}

// Adopted from https://github.com/moby/moby/blob/master/daemon/logger/ring.go#L128
// as this struct is not exported.
type ringBuffer struct {
	// A mutex lock is used here when writing/reading log messages from the queue
	// as there exists three go routines accessing the buffer.
	lock sync.Mutex
	// A condition variable wait is used here to notify goroutines that get access to
	// the buffer should wait or continue.
	wait *sync.Cond
	// current total bytes stored in the buffer
	curSizeInBytes int
	// maximum bytes capacity provided by the buffer
	maxSizeInBytes int
	// queue saves all the log messages read from pipes exposed by containerd, and
	// is consumed by underlying log driver.
	queue []*dockerlogger.Message
	// closedPipesCount is the number of closed container pipes for a single container.
	closedPipesCount int
	// isClosed indicates if ring buffer is closed.
	isClosed bool
}

// NewBufferedLogger creates a logger with the provided LoggerOpt,
// a buffer with customized max size and a channel monitor if stdout
// and stderr pipes are closed.
func NewBufferedLogger(l LogDriver, maxBufferSize int, containerID string) LogDriver {
	return &bufferedLogger{
		l:           l,
		buffer:      newLoggerBuffer(maxBufferSize),
		containerID: containerID,
	}
}

// newLoggerBuffer creates a buffer that stores messages which are
// from container and consumed by sub-level log drivers.
func newLoggerBuffer(maxBufferSize int) *ringBuffer {
	rb := &ringBuffer{
		maxSizeInBytes:   maxBufferSize,
		queue:            make([]*dockerlogger.Message, 0, ringCap),
		closedPipesCount: 0,
		isClosed:         false,
	}
	rb.wait = sync.NewCond(&rb.lock)

	return rb
}

// Start starts the non-blocking mode logger.
func (bl *bufferedLogger) Start(
	ctx context.Context,
	uid int,
	gid int,
	cleanupTime *time.Duration,
	ready func() error,
) error {
	pipeNameToPipe, err := bl.l.GetPipes()
	if err != nil {
		return err
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	// Start the goroutine of underlying log driver to consume logs from ring buffer and
	// send logs to destination when there's any.
	errGroup.Go(func() error {
		debug.SendEventsToLog(DaemonName, "Starting consuming logs from ring buffer", debug.INFO, 0)
		return bl.sendLogMessagesToDestination(uid, gid, cleanupTime)
	})

	// Start reading logs from container pipes.
	for pn, p := range pipeNameToPipe {
		// Copy pn and p to new variables source and pipe, accordingly.
		source := pn
		pipe := p

		errGroup.Go(func() error {
			logErr := bl.saveLogMessagesToRingBuffer(ctx, pipe, source, uid, gid)
			if logErr != nil {
				err := errors.Wrapf(logErr, "failed to send logs from pipe %s", source)
				debug.SendEventsToLog(DaemonName, err.Error(), debug.ERROR, 1)
				return err
			}
			return nil
		})
	}

	// Signal that the container is ready to be started
	if err := ready(); err != nil {
		return errors.Wrap(err, "failed to check container ready status")
	}

	// Wait() will return the first error it receives.
	return errGroup.Wait()
}

// saveLogMessagesToRingBuffer saves container log messages to ring buffer.
func (bl *bufferedLogger) saveLogMessagesToRingBuffer(
	ctx context.Context,
	f io.Reader,
	source string,
	uid int, gid int,
) error {
	if err := bl.Read(ctx, f, source, defaultBufSizeInBytes, bl.saveSingleLogMessageToRingBuffer); err != nil {
		err := errors.Wrapf(err, "failed to read logs from %s pipe", source)
		debug.SendEventsToLog(DaemonName, err.Error(), debug.ERROR, 1)
		return err
	}

	// No messages in the pipe, send signal to closed pipe channel.
	debug.SendEventsToLog(DaemonName, fmt.Sprintf("Pipe %s is closed", source), debug.INFO, 1)
	bl.buffer.closedPipesCount++
	// If both container pipes are closed, wake up the Dequeue goroutine which is waiting on wait.
	if bl.buffer.closedPipesCount == expectedNumOfPipes {
		bl.buffer.isClosed = true
		bl.buffer.wait.Broadcast()
	}

	return nil
}

// Read reads log messages from container pipe and saves them to ring buffer line by line.
func (bl *bufferedLogger) Read(
	ctx context.Context,
	pipe io.Reader,
	source string,
	bufferSizeInBytes int,
	sendLogMsgToDest sendLogToDestFunc,
) error {
	return bl.l.Read(ctx, pipe, source, bufferSizeInBytes, sendLogMsgToDest)
}

// saveSingleLogMessageToRingBuffer enqueues a single line of log message to ring buffer.
func (bl *bufferedLogger) saveSingleLogMessageToRingBuffer(
	line []byte,
	source string,
	isFirstPartial, isPartialMsg bool,
	partialTimestamp time.Time,
) (error, time.Time, bool, bool) {
	msgTimestamp, partialTimestamp, isFirstPartial, isPartialMsg := getLogTimestamp(
		isFirstPartial,
		isPartialMsg,
		partialTimestamp,
		bl.containerID,
	)
	if debug.Verbose {
		debug.SendEventsToLog(bl.containerID,
			fmt.Sprintf("[Pipe %s] Scanned message: %s", source, string(line)),
			debug.DEBUG, 0)
	}

	message := newMessage(line, source, msgTimestamp)
	err := bl.buffer.Enqueue(message)
	if err != nil {
		err := errors.Wrap(err, "failed to save logs to buffer")
		return err, partialTimestamp, isFirstPartial, isPartialMsg
	}

	return nil, partialTimestamp, isFirstPartial, isPartialMsg
}

// sendLogMessagesToDestination consumes logs from ring buffer and use the
// underlying log driver to send logs to destination.
func (bl *bufferedLogger) sendLogMessagesToDestination(uid int, gid int, cleanupTime *time.Duration) error {
	// Keep sending log message to destination defined by the underlying log driver until
	// the ring buffer is closed.
	for !bl.buffer.isClosed {
		if err := bl.sendLogMessageToDestination(); err != nil {
			debug.SendEventsToLog(DaemonName, err.Error(), debug.ERROR, 1)
			return err
		}
	}
	// If both container pipes are closed, flush messages left in ring buffer.
	debug.SendEventsToLog(DaemonName, "All pipes are closed, flushing buffer.", debug.INFO, 0)
	if err := bl.flushMessages(); err != nil {
		debug.SendEventsToLog(DaemonName, err.Error(), debug.ERROR, 1)
		return err
	}

	// Sleep sometime to let shim logger clean up, for example, to allow enough time for the last
	// few log messages be flushed to destination like CloudWatch.
	debug.SendEventsToLog(DaemonName,
		fmt.Sprintf("Sleeping %s for cleanning up.", cleanupTime.String()),
		debug.INFO, 0)
	time.Sleep(*cleanupTime)
	return nil
}

// sendLogMessageToDestination dequeues a single log message from buffer and sends to destination.
func (bl *bufferedLogger) sendLogMessageToDestination() error {
	msg, err := bl.buffer.Dequeue()
	// Do an early return if ring buffer is closed.
	if bl.buffer.isClosed {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to read logs from buffer")
	}

	err = bl.Log(msg.Line, msg.Source, msg.Timestamp)
	if err != nil {
		return errors.Wrap(err, "failed to send logs to destination")
	}

	return nil
}

// flushMessages flushes all the messages left in the ring buffer to
// destination after container pipes are closed.
func (bl *bufferedLogger) flushMessages() error {
	messages := bl.buffer.Flush()
	for _, msg := range messages {
		err := bl.Log(msg.Line, msg.Source, msg.Timestamp)
		if err != nil {
			return errors.Wrap(err, "unable to flush the remaining messages to destination")
		}
	}

	return nil
}

// Log lets underlying log driver send logs to destination.
func (bl *bufferedLogger) Log(line []byte, source string, logTimestamp time.Time) error {
	if debug.Verbose {
		debug.SendEventsToLog(DaemonName,
			fmt.Sprintf("[BUFFER] Sending message: %s", string(line)),
			debug.DEBUG, 0)
	}
	return bl.l.Log(line, source, logTimestamp)
}

// GetPipes gets pipes of container and its name that exposed by containerd.
func (bl *bufferedLogger) GetPipes() (map[string]io.Reader, error) {
	return bl.l.GetPipes()
}

// Adopted from https://github.com/moby/moby/blob/master/daemon/logger/ring.go#L155
// as messageRing struct is not exported.
// Enqueue adds a single log message to the tail of intermediate buffer.
func (b *ringBuffer) Enqueue(msg *dockerlogger.Message) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	lineSizeInBytes := len(msg.Line)
	// If there is already at least one log message in the queue and not enough space left
	// for the new coming log message to take up, drop this log message. Otherwise, save this
	// message to ring buffer anyway.
	if len(b.queue) > 0 &&
		b.curSizeInBytes+lineSizeInBytes > b.maxSizeInBytes {
		if debug.Verbose {
			debug.SendEventsToLog(DaemonName,
				"buffer is full/message is too long, waiting for available bytes",
				debug.DEBUG, 0)
			debug.SendEventsToLog(DaemonName,
				fmt.Sprintf("message size: %d, current buffer size: %d, max buffer size %d",
					lineSizeInBytes,
					b.curSizeInBytes,
					b.maxSizeInBytes),
				debug.DEBUG, 0)
		}

		// Wake up "Dequeue" or the other "Enqueue" go routine (called by the other pipe)
		// waiting on current mutex lock if there's any
		b.wait.Signal()
		return nil
	}

	b.queue = append(b.queue, msg)
	b.curSizeInBytes += lineSizeInBytes
	// Wake up "Dequeue" or the other "Enqueue" go routine (called by the other pipe)
	// waiting on current mutex lock if there's any
	b.wait.Signal()

	return nil
}

// Adopted from https://github.com/moby/moby/blob/master/daemon/logger/ring.go#L179
// as messageRing struct is not exported.
// Dequeue gets a line of log message from the head of intermediate buffer.
func (b *ringBuffer) Dequeue() (*dockerlogger.Message, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// If there is no log yet in the buffer, and the ring buffer is still open, wait
	// suspends current go routine.
	for len(b.queue) == 0 && !b.isClosed {
		if debug.Verbose {
			debug.SendEventsToLog(DaemonName,
				"No messages in queue, waiting...",
				debug.DEBUG, 0)
		}
		b.wait.Wait()
	}

	// Directly return if ring buffer is closed.
	if b.isClosed {
		return nil, nil
	}

	// Get and remove the oldest message saved in buffer/queue from head and update
	// the current used bytes of buffer.
	msg := b.queue[0]
	b.queue = b.queue[1:]
	b.curSizeInBytes -= len(msg.Line)

	return msg, nil
}

// Adopted from https://github.com/moby/moby/blob/master/daemon/logger/ring.go#L215
// as messageRing struct is not exported.
// Flush flushes all the messages left in the buffer and clear queue.
func (b *ringBuffer) Flush() []*dockerlogger.Message {
	b.lock.Lock()
	defer b.lock.Unlock()

	if len(b.queue) == 0 {
		return make([]*dockerlogger.Message, 0)
	}

	messages := b.queue
	b.queue = make([]*dockerlogger.Message, 0)

	return messages
}
