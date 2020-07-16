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
	"fmt"
	"io"
	"math"
	"syscall"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"

	"github.com/coreos/go-systemd/journal"
	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/pkg/errors"
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

	// defaultBufSizeInBytes provides a reasonable default for loggers that do
	// not have an external limit to impose on log line size. Adopted this value
	// from Docker, reference:
	// https://github.com/moby/moby/blob/19.03/daemon/logger/copier.go#L21
	defaultBufSizeInBytes = 16 * 1024
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
	GetPipes() (io.Reader, io.Reader)
	// Log sends logs to destination.
	Log([]byte, string, time.Time) error
}

// NewLogger creates a LogDriver with the provided LoggerOpt
func NewLogger(options ...LoggerOpt) (LogDriver, error) {
	l := &Logger{
		Info:              &dockerlogger.Info{},
		bufferSizeInBytes: defaultBufSizeInBytes,
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
	if l.Stdout == nil || l.Stderr == nil {
		return errors.New("no stdout/stderr pipe opened")
	}
	pipeNameToPipe := map[string]io.Reader{
		sourceSTDOUT: l.Stdout,
		sourceSTDERR: l.Stderr,
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	for pn, p := range pipeNameToPipe {
		// Copy pn and p to new variables source and pipe, accordingly.
		source := pn
		pipe := p

		errGroup.Go(func() error {
			logErr := l.sendLogs(ctx, pipe, source, uid, gid, cleanupTime)
			if logErr != nil {
				err := errors.Wrapf(logErr, "failed to send logs from pipe %s", source)
				debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr, 1)
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

// sendLogs sends logs to destination.
func (l *Logger) sendLogs(
	ctx context.Context,
	f io.Reader,
	source string,
	uid int, gid int,
	cleanupTime *time.Duration,
) error {
	// Set uid and/or gid for this goroutine. Currently the Setuid/SetGID syscall does not
	// apply on threads in golang, see issue: https://github.com/golang/go/issues/1435
	// TODO: remove it once the changes are released: https://go-review.googlesource.com/c/go/+/210639
	if err := SetUIDAndGID(uid, gid); err != nil {
		debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr, 1)
		return err
	}

	if err := l.read(ctx, f, source); err != nil {
		err := errors.Wrapf(err, "failed to read logs from %s pipe", source)
		debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr, 1)
		return err
	}

	// Sleep sometime to let shim logger clean up, for example, to allow enough time for the last
	// few log messages be flushed to destination like CloudWatch.
	debug.SendEventsToJournal(DaemonName,
		fmt.Sprintf("Pipe %s is closed. Sleeping %s for cleanning up.", source, cleanupTime.String()),
		journal.PriInfo,
		0)
	time.Sleep(*cleanupTime)
	return nil
}

// read gets container logs, saves them to our own buffer. Then we will read logs line by line
// and send them to destination. More log messages will be sent to system journal in verbose mdoe
// for debugging.
func (l *Logger) read(ctx context.Context, pipe io.Reader, source string) error {
	var (
		partialTimestamp time.Time
		bytesInBuffer    int
		err              error
		eof              bool
	)
	// Initiate an in-memory buffer to hold bytes read from container pipe.
	buf := make([]byte, l.bufferSizeInBytes)
	// isFirstPartial indicates if current message saved in buffer is not a complete line,
	// and is the first partial of the whole log message. Initialize to true.
	isFirstPartial := true
	// isPartialMsg indicates if current message is a partial log message. Initialize to false.
	isPartialMsg := false
	for {
		select {
		case <-ctx.Done():
			debug.SendEventsToJournal(l.Info.ContainerID,
				fmt.Sprintf("Logging stopped in pipe %s", source),
				journal.PriDebug, 0)
			return nil
		default:
			eof, bytesInBuffer, err = l.readFromContainerPipe(pipe, buf, bytesInBuffer)
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
				curLine := buf[head : head+lenOfLine]
				err, partialTimestamp, _, _ = l.sendLogMsgToDest(
					curLine,
					source,
					isFirstPartial,
					isPartialMsg,
					partialTimestamp,
				)
				if err != nil {
					return err
				}
				// Since we have found a newline symbol, it means this line has ended.
				// Reset flags.
				isFirstPartial = true
				isPartialMsg = false

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
					err, partialTimestamp, isFirstPartial, isPartialMsg = l.sendLogMsgToDest(
						curLine,
						source,
						isFirstPartial,
						true, // Record as a partial message.
						partialTimestamp,
					)
					if err != nil {
						return err
					}

					// reset head and bytesInBuffer
					head = 0
					bytesInBuffer = 0
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

// readFromContainerPipe reads bytes from container pipe, upto max read size in bytes of 2048.
func (l *Logger) readFromContainerPipe(pipe io.Reader, buf []byte, bytesInBuffer int) (bool, int, error) {
	// eof indicates if we have already met EOF error.
	eof := false
	// Decide how many bytes we can read from container pipe for this iteration. It's either
	// the current bytes in buffer plus 2048 bytes or the available spaces left, whichever is
	// smaller.
	readBytesUpto := int(math.Min(float64(bytesInBuffer+l.maxReadBytes), float64(cap(buf))))
	// Read logs from container pipe if there are available spaces for new log messages.
	if readBytesUpto > bytesInBuffer {
		readBytesFromPipe, err := pipe.Read(buf[bytesInBuffer:readBytesUpto])
		if err != nil {
			if err != io.EOF {
				return false, bytesInBuffer, errors.Wrap(err, "failed to read log stream from container pipe")
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
	isFirstPartial, isPartialMsg bool,
	partialTimestamp time.Time,
) (error, time.Time, bool, bool) {
	msgTimestamp, partialTimestamp, isFirstPartial, isPartialMsg := l.getLogTimestamp(
		isFirstPartial,
		isPartialMsg,
		partialTimestamp,
	)
	if debug.Verbose {
		debug.SendEventsToJournal(l.Info.ContainerID,
			fmt.Sprintf("[Pipe %s] Scanned message: %s", source, string(line)),
			journal.PriDebug, 0)
	}

	err := l.Log(line, source, msgTimestamp)
	if err != nil {
		return err, partialTimestamp, isFirstPartial, isPartialMsg
	}

	return nil, partialTimestamp, isFirstPartial, isPartialMsg
}

// getLogTimestamp gets the timestamp of a log message. It could be current timestamp
// if it's a new line, or the recorded timestamp from the first partial if it's the other partial
// messages.
func (l *Logger) getLogTimestamp(
	isFirstPartial, isPartialMsg bool,
	partialTimestamp time.Time,
) (time.Time, time.Time, bool, bool) {
	msgTimestamp := time.Now().UTC()

	// If it is not end with with newline for the first time, record it as a partial
	// message and set it to be the first partial. Then record the timestamp as the
	// timestamp of the whole log message.
	if isFirstPartial {
		partialTimestamp = msgTimestamp
		if debug.Verbose {
			debug.SendEventsToJournal(l.Info.ContainerID,
				fmt.Sprintf("Saving first partial at time %s", partialTimestamp.String()),
				journal.PriDebug, 0)
		}
		// Set isFirstPartial to false and set indicator of partial log message to be true.
		isFirstPartial = false
		isPartialMsg = true
	} else if isPartialMsg {
		// If there are more partial messages recorded before the current read, use the
		// recorded timestamp as it of the current message as well.
		msgTimestamp = partialTimestamp
		if debug.Verbose {
			debug.SendEventsToJournal(l.Info.ContainerID,
				fmt.Sprintf("Setting partial log message to time %s", msgTimestamp.String()),
				journal.PriDebug, 0)
		}
	}

	return msgTimestamp, partialTimestamp, isFirstPartial, isPartialMsg
}

// Log sends logs to destination.
func (l *Logger) Log(line []byte, source string, logTimestamp time.Time) error {
	message := newMessage(line, source, logTimestamp)
	err := l.Stream.Log(message)
	if err != nil {
		return errors.Wrapf(err, "failed to log msg for container %s", l.Info.ContainerName)
	}

	return nil
}

// newMessage creates a new logger message.
func newMessage(line []byte, source string, logTimestamp time.Time) *dockerlogger.Message {
	msg := dockerlogger.NewMessage()
	msg.Line = append(msg.Line, line...)
	msg.Source = source
	msg.Timestamp = logTimestamp

	return msg
}

// GetPipes gets pipes of container that exposed by containerd.
func (l *Logger) GetPipes() (io.Reader, io.Reader) {
	return l.Stdout, l.Stderr
}

// SetUIDAndGID sets UID and/or GID for current goroutine.
// TODO: move it to main package once the changes are released: https://go-review.googlesource.com/c/go/+/210639
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

// setUID sets UID of current goroutine.
func setUID(id int) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_SETUID, uintptr(id), 0, 0); errno != 0 {
		return errors.Wrap(errors.New(errno.Error()), "unable to set uid")
	}

	// Check if uid set correctly
	u := syscall.Getuid()
	if u != id {
		return errors.New(fmt.Sprintf("want uid %d, but get uid %d", id, u))
	}
	debug.SendEventsToJournal(DaemonName,
		fmt.Sprintf("Set uid: %d", u),
		journal.PriInfo, 1)

	return nil
}

// setGID sets GID of current goroutine.
func setGID(id int) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_SETGID, uintptr(id), 0, 0); errno != 0 {
		return errors.Wrap(errors.New(errno.Error()), "unable to set gid")
	}

	// Check if gid set correctly
	g := syscall.Getgid()
	if g != id {
		return errors.New(fmt.Sprintf("want gid %d, but get gid %d", id, g))
	}
	debug.SendEventsToJournal(DaemonName,
		fmt.Sprintf("Set gid %d", g),
		journal.PriInfo, 1)

	return nil
}
