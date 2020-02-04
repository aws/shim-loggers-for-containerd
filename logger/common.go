package logger

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"syscall"
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
	UID           int
	GID           int
}

// Basic Logger struct for all log drivers
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
	Start(int, int, func() error) error
	// GetPipes gets pipes of container that exposed by containerd.
	GetPipes() (io.Reader, io.Reader)
	// LogWithRetry sends logs to destination with retry.
	LogWithRetry(line []byte, logTimestamp time.Time) error
}

// NewLogger creates a LogDriver with the provided LoggerOpt
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

// Start starts the actual logger.
func (l *Logger) Start(uid int, gid int, ready func() error) error {
	var wg sync.WaitGroup
	if l.Stdout != nil {
		wg.Add(1)
		go l.sendLogs(l.Stdout, &wg, uid, gid)
	}
	if l.Stderr != nil {
		wg.Add(1)
		go l.sendLogs(l.Stderr, &wg, uid, gid)
	}

	// Signal that the container is ready to be started
	if err := ready(); err != nil {
		return errors.Wrap(err, "failed to check container ready status")
	}
	wg.Wait()

	return nil
}

// sendLogs sends logs to destintion.
func (l *Logger) sendLogs(f io.Reader, wg *sync.WaitGroup, uid int, gid int) {
	defer wg.Done()

	// Set uid and/or gid for this goroutine. Currently the Setuid/SetGID syscall does not
	// apply on threads in golang, see issue: https://github.com/golang/go/issues/1435
	// TODO: remove it once the changes are released: https://go-review.googlesource.com/c/go/+/210639
	if err := SetUIDAndGID(uid, gid); err != nil {
		debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr)
		return
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			debug.SendEventsToJournal(DaemonName,
				"Message is empty, skip saving", journal.PriInfo)
			continue
		}
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

// GetPipes gets pipes of container that exposed by containerd.
func (l *Logger) GetPipes() (io.Reader, io.Reader) {
	return l.Stdout, l.Stderr
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
		journal.PriInfo)

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
		journal.PriInfo)

	return nil
}
