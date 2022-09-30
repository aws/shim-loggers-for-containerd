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
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"

	dockerlogger "github.com/docker/docker/daemon/logger"

	types "github.com/docker/docker/api/types/backend"
	"golang.org/x/sync/errgroup"
)

const (
	DaemonName      = "shim-loggers-for-containerd"
	NonBlockingMode = "non-blocking"

	// source pipe of log message
	sourceSTDOUT = "stdout"
	sourceSTDERR = "stderr"

	newline = '\n'

	// defaultMaxReadBytes is the default maximum bytes read during a single read
	// operation. Adopted this value from Docker, reference:
	// https://github.com/moby/moby/blob/19.03/daemon/logger/copier.go#L17
	defaultMaxReadBytes = 2 * 1024

	// DefaultBufSizeInBytes provides a reasonable default for loggers that do
	// not have an external limit to impose on log line size. Adopted this value
	// from Docker, reference:
	// https://github.com/moby/moby/blob/19.03/daemon/logger/copier.go#L21
	DefaultBufSizeInBytes = 16 * 1024
)

type GlobalArgs struct {
	// Required arguments
	ContainerID   string
	ContainerName string
	LogDriver     string

	// Optional arguments
	Mode          string
	MaxBufferSize int
	UID           int
	GID           int
	CleanupTime   *time.Duration
}

// Optional docker config arguments
type DockerConfigs struct {
	ContainerImageID   string
	ContainerImageName string
	ContainerEnv       []string
	ContainerLabels    map[string]string
}

// Basic Logger struct for all log drivers
type Logger struct {
	Info   *dockerlogger.Info
	Stream Client
	Stdout io.Reader
	Stderr io.Reader

	// bufferSizeInBytes defines the size of our own buffer. It's default to
	// 16 * 1024, but maybe different among log drivers.
	bufferSizeInBytes int
	// maxReadBytes defines how many bytes we want to read from container pipe
	// per iteration. It's default to 2 * 1024.
	maxReadBytes int
}

// WindowsArgs struct for Windows configuration
type WindowsArgs struct {
	ProxyEnvVar string
	LogFileDir  string
}

// Client is a wrapper for docker logger's Log method, which is mostly used for testing
// purposes.
type Client interface {
	Log(*dockerlogger.Message) error
}

// Interface for all log drivers
type LogDriver interface {
	// Start functions starts sending container logs to destination.
	Start(context.Context, int, int, *time.Duration, func() error) error
	// GetPipes gets pipes of container that exposed by containerd.
	GetPipes() (map[string]io.Reader, error)
	// Log sends logs to destination.
	Log(*dockerlogger.Message) error
	// Read reads a single log message from container pipe and sends it to
	// destination or saves it to ring buffer, depending on the mode of log
	// driver.
	Read(context.Context, io.Reader, string, int, sendLogToDestFunc) error
}

// NewLogger creates a LogDriver with the provided LoggerOpt
func NewLogger(options ...LoggerOpt) (LogDriver, error) {
	l := &Logger{
		Info:              &dockerlogger.Info{},
		bufferSizeInBytes: DefaultBufSizeInBytes,
		maxReadBytes:      defaultMaxReadBytes,
	}
	for _, opt := range options {
		opt(l)
	}
	return l, nil
}

// Placeholder info. Expected that relevant parts will be modified
// via the common_opts.
func NewInfo(containerID string, containerName string, options ...InfoOpt) *dockerlogger.Info {
	info := &dockerlogger.Info{
		Config:           make(map[string]string),
		ContainerID:      containerID,
		ContainerName:    containerName,
		ContainerArgs:    make([]string, 0),
		ContainerCreated: time.Now(),
		ContainerEnv:     make([]string, 0),
		ContainerLabels:  make(map[string]string),
		DaemonName:       DaemonName,
	}

	for _, opt := range options {
		opt(info)
	}

	return info
}

// UpdateDockerConfigs updates the docker config fields to the logger info.
func UpdateDockerConfigs(info *dockerlogger.Info, dockerConfigs *DockerConfigs) *dockerlogger.Info {
	info.ContainerImageName = dockerConfigs.ContainerImageName
	info.ContainerImageID = dockerConfigs.ContainerImageID
	info.ContainerLabels = dockerConfigs.ContainerLabels
	info.ContainerEnv = dockerConfigs.ContainerEnv
	return info
}

// Start starts the actual logger.
func (l *Logger) Start(
	ctx context.Context,
	uid int,
	gid int,
	cleanupTime *time.Duration,
	ready func() error,
) error {
	pipeNameToPipe, err := l.GetPipes()
	if err != nil {
		return err
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	for pn, p := range pipeNameToPipe {
		// Copy pn and p to new variables source and pipe, accordingly.
		source := pn
		pipe := p

		errGroup.Go(func() error {
			logErr := l.sendLogs(ctx, pipe, source, uid, gid, cleanupTime)
			if logErr != nil {
				err := fmt.Errorf("failed to send logs from pipe %s: %w", source, err)
				debug.SendEventsToLog(DaemonName, err.Error(), debug.ERROR, 1)
				return err
			}
			return nil
		})
	}

	// Signal that the container is ready to be started
	if err := ready(); err != nil {
		return fmt.Errorf("failed to check container ready status: %w", err)
	}

	// Wait() will return the first error it receives.
	return errGroup.Wait()
}

// sendLogs sends logs to destination.
func (l *Logger) sendLogs(
	ctx context.Context,
	f io.Reader,
	source string,
	uid int, gid int,
	cleanupTime *time.Duration,
) error {
	if err := l.Read(ctx, f, source, l.bufferSizeInBytes, l.sendLogMsgToDest); err != nil {
		err := fmt.Errorf("failed to read logs from %s pipe: %w", source, err)
		debug.SendEventsToLog(DaemonName, err.Error(), debug.ERROR, 1)
		return err
	}

	// Sleep sometime to let shim logger clean up, for example, to allow enough time for the last
	// few log messages be flushed to destination like CloudWatch.
	debug.SendEventsToLog(DaemonName,
		fmt.Sprintf("Pipe %s is closed. Sleeping %s for cleanning up.", source, cleanupTime.String()),
		debug.INFO,
		0)
	time.Sleep(*cleanupTime)
	return nil
}

// sendLogToDestFunc is type a function that gets used in read function, which is defined by
// the underlying logger.
type sendLogToDestFunc func(
	line []byte,
	source string,
	isPartialMsg, isLastPartial bool,
	partialID string,
	partialOrdinal int,
	msgTimestamp time.Time,
) error

// Read gets container logs, saves them to our own buffer. Then we will read logs line by line
// and send them to destination. In non-blocking mode, the destination is the ring buffer. More
// log messages will be sent in verbose mode for debugging.
func (l *Logger) Read(
	ctx context.Context,
	pipe io.Reader,
	source string,
	bufferSizeInBytes int,
	sendLogMsgToDest sendLogToDestFunc,
) error {
	var (
		msgTimestamp  time.Time
		bytesInBuffer int
		err           error
		eof           bool
	)
	// Initiate an in-memory buffer to hold bytes read from container pipe.
	buf := make([]byte, bufferSizeInBytes)
	// isFirstPartial indicates if current message saved in buffer is not a complete line,
	// and is the first partial of the whole log message. Initialize to true.
	isFirstPartial := true
	// isPartialMsg indicates if current message is a partial log message. Initialize to false.
	isPartialMsg := false
	// isLastPartial indicates if this message completes a partial message
	isLastPartial := false
	// partialID is a random ID given to each split message
	partialID := ""
	// partialOrdinal orders the split messages and count up from 1
	partialOrdinal := 1

	for {
		select {
		case <-ctx.Done():
			debug.SendEventsToLog(l.Info.ContainerID,
				fmt.Sprintf("Logging stopped in pipe %s", source),
				debug.DEBUG, 0)
			return nil
		default:
			eof, bytesInBuffer, err = readFromContainerPipe(pipe, buf, bytesInBuffer, l.maxReadBytes)
			if err != nil {
				return err
			}
			// If container pipe is closed and no bytes left in our buffer, directly return.
			if eof && bytesInBuffer == 0 {
				return nil
			}

			// Iteratively scan the unread part in our own buffer and read logs line by line.
			// Then send it to destination.
			head := 0
			// This function returns -1 if '\n' in not present in buffer.
			lenOfLine := bytes.IndexByte(buf[head:bytesInBuffer], newline)
			for lenOfLine >= 0 {
				// If this is the end of a partial message
				// use the existing timestamp, so that all
				// partials split from the same message have the same timestamp
				// If not, new timestamp.
				if isPartialMsg {
					isLastPartial = true
				} else {
					msgTimestamp = time.Now().UTC()
				}
				curLine := buf[head : head+lenOfLine]
				err = sendLogMsgToDest(
					curLine,
					source,
					isPartialMsg,
					isLastPartial,
					partialID,
					partialOrdinal,
					msgTimestamp,
				)
				if err != nil {
					return err
				}
				// Since we have found a newline symbol, it means this line has ended.
				// Reset flags.
				isFirstPartial = true
				isPartialMsg = false
				isLastPartial = false
				partialID = ""
				partialOrdinal = 1

				// Update the index of head of next line message.
				head += lenOfLine + 1
				lenOfLine = bytes.IndexByte(buf[head:bytesInBuffer], newline)
			}

			// If the pipe is closed and the last line does not end with a newline symbol, send whatever left
			// in the buffer to destination as a single log message. Or if our buffer is full but there is
			// no newline symbol yet, record it as a partial log message and send it as a single log message
			// to destination.
			if eof || bufferIsFull(buf, head, bytesInBuffer) {
				// Still bytes left in the buffer after we identified all newline symbols.
				if head < bytesInBuffer {
					curLine := buf[head:bytesInBuffer]

					// Record as a partial message.
					isPartialMsg = true
					if isFirstPartial {
						msgTimestamp = time.Now().UTC()
						partialID, err = generateRandomID()
					}
					if err != nil {
						return err
					}

					err = sendLogMsgToDest(
						curLine,
						source,
						isPartialMsg,
						isLastPartial,
						partialID,
						partialOrdinal,
						msgTimestamp,
					)
					if err != nil {
						return err
					}

					// reset head and bytesInBuffer
					head = 0
					bytesInBuffer = 0
					// increment partial flags
					partialOrdinal += 1
					if isFirstPartial {
						// if this was the first partial message
						// the next one is not the first if it is also partial
						isFirstPartial = false
					}
				}
				// If pipe is closed after we send all bytes left in buffer, then directly return.
				if eof {
					return nil
				}
			}

			// If there are any bytes left in the buffer, move them to the head and handle them in the
			// next round.
			if head > 0 {
				copy(buf[0:], buf[head:bytesInBuffer])
				bytesInBuffer -= head
			}
		}
	}
}

// generateRandomID is based on Docker
// GenerateRandomID: https://github.com/moby/moby/blob/bca8d9f2ce0d63e1490692917cde6273bc288bad/pkg/stringid/stringid.go#L40
// with the simplification that we don't need to worry about guaranteeing the string isn't all 0 - 9
// Consequently ^ we have our own function instead of importing from Docker.
func generateRandomID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	id := hex.EncodeToString(b)
	return id, nil
}

// readFromContainerPipe reads bytes from container pipe, upto max read size in bytes of 2048.
func readFromContainerPipe(pipe io.Reader, buf []byte, bytesInBuffer, maxReadBytes int) (bool, int, error) {
	// eof indicates if we have already met EOF error.
	eof := false
	// Decide how many bytes we can read from container pipe for this iteration. It's either
	// the current bytes in buffer plus 2048 bytes or the available spaces left, whichever is
	// smaller.
	readBytesUpto := int(math.Min(float64(bytesInBuffer+maxReadBytes), float64(cap(buf))))
	// Read logs from container pipe if there are available spaces for new log messages.
	if readBytesUpto > bytesInBuffer {
		readBytesFromPipe, err := pipe.Read(buf[bytesInBuffer:readBytesUpto])
		if err != nil {
			if err != io.EOF {
				return false, bytesInBuffer, fmt.Errorf("failed to read log stream from container pipe: %w", err)
			}
			// Pipe is closed, set flag to true.
			eof = true
		}
		bytesInBuffer += readBytesFromPipe
	}

	return eof, bytesInBuffer, nil
}

// bufferIsFull indicates if our own buffer is full.
func bufferIsFull(buf []byte, head, bytesInBuffer int) bool {
	return head == 0 && bytesInBuffer == len(buf)
}

// sendLogMsgToDest sends a single line of log message to destination.
func (l *Logger) sendLogMsgToDest(
	line []byte,
	source string,
	isPartialMsg, isLastPartial bool,
	partialID string,
	partialOrdinal int,
	msgTimestamp time.Time,
) error {
	if debug.Verbose {
		debug.SendEventsToLog(l.Info.ContainerID,
			fmt.Sprintf("[Pipe %s] Scanned message: %s", source, string(line)),
			debug.DEBUG, 0)
	}

	message := newMessage(line, source, msgTimestamp)
	if isPartialMsg {
		message.PLogMetaData = &types.PartialLogMetaData{ID: partialID, Ordinal: partialOrdinal, Last: isLastPartial}
	}
	err := l.Log(message)
	if err != nil {
		return fmt.Errorf("failed to log msg for container %s: %w", l.Info.ContainerName, err)
	}

	return nil
}

// Log sends logs to destination.
func (l *Logger) Log(message *dockerlogger.Message) error {
	return l.Stream.Log(message)
}

// newMessage creates a new logger message.
func newMessage(line []byte, source string, logTimestamp time.Time) *dockerlogger.Message {
	msg := dockerlogger.NewMessage()
	msg.Line = append(msg.Line, line...)
	msg.Source = source
	msg.Timestamp = logTimestamp

	return msg
}

// GetPipes gets pipes of container and its name that exposed by containerd.
func (l *Logger) GetPipes() (map[string]io.Reader, error) {
	if l.Stdout == nil || l.Stderr == nil {
		return nil, errors.New("no stdout/stderr pipe opened")
	}

	pipeNameToPipe := map[string]io.Reader{
		sourceSTDOUT: l.Stdout,
		sourceSTDERR: l.Stderr,
	}

	return pipeNameToPipe, nil
}

// SetUIDAndGID sets UID and/or GID for current goroutine/process.
// If you are building with go version includes the following commit, you only need to call this once. Otherwise
// you need call this function in all goroutines.
// Commit: https://github.com/golang/go/commit/d1b1145cace8b968307f9311ff611e4bb810710c
// TODO: remove the above comment once the changes are released: https://go-review.googlesource.com/c/go/+/210639
func SetUIDAndGID(uid int, gid int) error {
	// gid<0 is assumed as gid argument is not set and is directly ignored.
	switch {
	case gid == 0:
		// gid=0 is not supported in shim logger.
		return errors.New("setting gid with value of zero is not supported")
	case gid > 0:
		if err := setGID(gid); err != nil {
			return err
		}
	}

	// uid<0 is assumed as uid argument is not set and is directly ignored.
	switch {
	case uid == 0:
		// uid=0 is not supported in shim logger.
		return errors.New("setting uid with value of zero is not supported")
	case uid > 0:
		if err := setUID(uid); err != nil {
			return err
		}
	}

	return nil
}
