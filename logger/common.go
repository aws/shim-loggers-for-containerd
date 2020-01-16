package logger

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"

	"github.com/aws/amazon-ecs-agent/agent/utils/retry"
	"github.com/coreos/go-systemd/journal"
	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/pkg/errors"
)

const (
	DaemonName      = "shim-loggers-for-containerd"
	NonBlockingMode = "non-blocking"

	// Define the retry parameters for retrying driving logs to destination
	LogRetryMaxAttempts = 3
	LogRetryMinBackoff  = 500 * time.Millisecond
	LogRetryMaxBackoff  = 1 * time.Second
	LogRetryJitter      = 0.3
	LogRetryMultiple    = 2
)

type GlobalArgs struct {
	// Required arguments
	ContainerID   string
	ContainerName string
	LogDriver     string

	// Optional arguments
	Mode          string
	MaxBufferSize int
}

type Logger struct {
	Info   *dockerlogger.Info
	Stream client
	Stdout io.Reader
	Stderr io.Reader
}

// client is a wrapper for docker logger's Log method, which is mostly used for testing
// purposes.
type client interface {
	Log(*dockerlogger.Message) error
}

// Interface for all log drivers
type LogDriver interface {
	// Start functions starts sending container logs to destination.
	Start(func() error) error
}

// NewLogger creates a logDriver with the provided LoggerOpt
func NewLogger(options ...LoggerOpt) (LogDriver, error) {
	l := &Logger{
		Info: &dockerlogger.Info{},
	}
	for _, opt := range options {
		opt(l)
	}
	return l, nil
}

// Placeholder info. Expected that relevant parts will be modified
// via the logger_opts.
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

// Start starts the actual logger.
func (l *Logger) Start(ready func() error) error {
	var wg sync.WaitGroup
	if l.Stdout != nil {
		wg.Add(1)
		go l.sendLogs(l.Stdout, &wg)
	}
	if l.Stderr != nil {
		wg.Add(1)
		go l.sendLogs(l.Stderr, &wg)
	}

	// Signal that the container is ready to be started
	if err := ready(); err != nil {
		return errors.Wrap(err, "failed to check container ready status")
	}
	wg.Wait()

	return nil
}

// sendLogs sends logs to destintion.
func (l *Logger) sendLogs(f io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if err := l.read(scanner); err != nil {
			debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr)
			return
		}
	}
}

// read gets container logs and sends to destination.
// More log messages will be sent to system journal
// in verbose mdoe for debugging.
func (l *Logger) read(s *bufio.Scanner) error {
	if s.Err() != nil {
		return errors.Wrap(s.Err(), "failed to get logs from container")
	}

	if debug.Verbose {
		debug.SendEventsToJournal(l.Info.ContainerID,
			fmt.Sprintf("[SCANNER] Scanned message: %s", s.Text()),
			journal.PriDebug)
	}
	// Send logs to destination with underlying log driver
	err := l.LogWithRetry(s.Bytes(), time.Now())
	if err != nil {
		return errors.Wrap(s.Err(), "failed to send logs to destination")
	}

	return nil
}

// LogWithRetry sends logs to destination with retry.
func (l *Logger) LogWithRetry(line []byte, logTimestamp time.Time) error {
	retryTimes := 0
	message := newMessage(line, l.Info.ContainerID, logTimestamp)
	backoff := newBackoff()
	err := retry.RetryNWithBackoff(
		backoff,
		LogRetryMaxAttempts,
		func() error {
			retryTimes += 1
			return l.Stream.Log(message)
		})
	if err != nil {
		err = errors.Wrapf(err, "sending container logs to destination has been retried for %d times", retryTimes)
		return err
	}

	return nil
}

// newBackoff creates a new Backoff object.
func newBackoff() retry.Backoff {
	return retry.NewExponentialBackoff(
		LogRetryMinBackoff,
		LogRetryMaxBackoff,
		LogRetryJitter,
		LogRetryMultiple)
}

// newMessage creates a new logger message.
func newMessage(line []byte, source string, logTimestamp time.Time) *dockerlogger.Message {
	msg := dockerlogger.NewMessage()
	msg.Line = line
	msg.Source = source
	msg.Timestamp = logTimestamp

	return msg
}
