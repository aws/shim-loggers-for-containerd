package awslogs

import (
	"bufio"
	"context"
	"io"
	"sync"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"

	"github.com/aws/amazon-ecs-agent/agent/utils/retry"
	"github.com/containerd/containerd/runtime/v2/logging"
	"github.com/coreos/go-systemd/journal"
	dockerlogger "github.com/docker/docker/daemon/logger"
	dockerawslogs "github.com/docker/docker/daemon/logger/awslogs"
	"github.com/pkg/errors"
)

const (
	nonBlockingMode = "non-blocking"

	// awslogs driver options
	RegionKey              = "awslogs-region"
	GroupKey               = "awslogs-group"
	CreateGroupKey         = "awslogs-create-group"
	StreamKey              = "awslogs-stream"
	MultilinePatternKey    = "awslogs-multiline-pattern"
	DatetimeFormatKey      = "awslogs-datetime-format"
	CredentialsEndpointKey = "awslogs-credentials-endpoint"
)

// AWS logs driver specified arguments
type Args struct {
	// Required arguments
	Group               string
	Region              string
	Stream              string
	CredentialsEndpoint string

	// Optional arguments
	CreateGroup      string
	MultilinePattern string
	DatetimeFormat   string
}

// LoggerArgs stores global logger args and awslogs specific args
type LoggerArgs struct {
	globalArgs *logger.GlobalArgs
	args       *Args
}

type logDriver struct {
	info   *dockerlogger.Info
	stream client

	stdout io.Reader
	stderr io.Reader
}

// InitLogger initialize the input arguments
func InitLogger(globalArgs *logger.GlobalArgs, awslogsArgs *Args) *LoggerArgs {
	return &LoggerArgs{
		globalArgs: globalArgs,
		args:       awslogsArgs,
	}
}

// client is a wrapper for docker logger's Log method, which is mostly used for testing
// purposes.
type client interface {
	Log(*dockerlogger.Message) error
}

// NewLogger creates a awslogs logDriver with the provided LoggerOpt
func NewLogger(options ...LoggerOpt) (logger.LogDriver, error) {
	l := &logDriver{
		info: &dockerlogger.Info{},
	}
	for _, opt := range options {
		opt(l)
	}
	stream, err := dockerawslogs.New(*l.info)
	if err != nil {
		err = errors.Wrap(err, "unable to create awslogs driver")
		return nil, err
	}
	l.stream = stream
	return l, nil
}

// RunLogDriver initiates an awslogs driver and starts driving container logs to cloudwatch
func (la *LoggerArgs) RunLogDriver(ctx context.Context, config *logging.Config, ready func() error) error {
	loggerConfig := getAWSLogsConfig(la.args)
	info := newInfo(
		la.globalArgs.ContainerID,
		la.globalArgs.ContainerName,
		WithConfig(loggerConfig),
	)
	l, err := NewLogger(
		WithStdout(config.Stdout),
		WithStderr(config.Stderr),
		WithInfo(info),
	)
	if err != nil {
		return err
	}

	if la.globalArgs.Mode == nonBlockingMode {
		l = logger.NewBufferedLogger(l, la.globalArgs.MaxBufferSize)
	}

	// Start awslogs driver
	debug.SendEventsToJournal(logger.DaemonName, "Starting awslogs driver", journal.PriInfo)
	err = l.Start(ready)
	if err != nil {
		return errors.Wrap(err, "failed to start awslogs driver")
	}

	return nil
}

// Placeholder info. Expected that relevant parts will be modified
// via the logger_opts.
func newInfo(containerID string, containerName string, options ...InfoOpt) *dockerlogger.Info {
	info := &dockerlogger.Info{
		Config:           make(map[string]string),
		ContainerID:      containerID,
		ContainerName:    containerName,
		ContainerArgs:    make([]string, 0),
		ContainerCreated: time.Now(),
		ContainerEnv:     make([]string, 0),
		ContainerLabels:  make(map[string]string),
		DaemonName:       logger.DaemonName,
	}

	for _, opt := range options {
		opt(info)
	}

	return info
}

// getAWSLogsConfig sets values for awslogs config
func getAWSLogsConfig(args *Args) map[string]string {
	config := make(map[string]string)
	// Required arguments
	config[GroupKey] = args.Group
	config[RegionKey] = args.Region
	config[StreamKey] = args.Stream
	config[CredentialsEndpointKey] = args.CredentialsEndpoint
	// Optional arguments
	createGroup := args.CreateGroup
	if createGroup != "" {
		config[CreateGroupKey] = createGroup
	}
	multilinePattern := args.MultilinePattern
	if multilinePattern != "" {
		config[MultilinePatternKey] = multilinePattern
	}
	datetimeFormat := args.DatetimeFormat
	if datetimeFormat != "" {
		config[DatetimeFormatKey] = datetimeFormat
	}

	return config
}

// Start starts the actual logger.
func (l *logDriver) Start(ready func() error) error {
	var wg sync.WaitGroup
	if l.stdout != nil {
		wg.Add(1)
		go l.sendLogs(l.stdout, &wg)
	}
	if l.stderr != nil {
		wg.Add(1)
		go l.sendLogs(l.stderr, &wg)
	}

	// Signal that the container is ready to be started
	if err := ready(); err != nil {
		return errors.Wrap(err, "failed to check container ready status")
	}
	wg.Wait()

	return nil
}

// sendLogs sends logs to aws cloudwatch logs.
func (l *logDriver) sendLogs(f io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if err := l.read(scanner); err != nil {
			debug.SendEventsToJournal(logger.DaemonName, err.Error(), journal.PriErr)
			return
		}
	}
}

// read gets container logs and sends to cloudwatch logs.
// For now we send the log messages to system journal
// as well for debugging. We may remove it in the future.
func (l *logDriver) read(s *bufio.Scanner) error {
	if s.Err() != nil {
		return errors.Wrap(s.Err(), "failed to get logs from container")
	}

	// Leave for debugging purpose.
	// Container logs are identified by it's container ID in system journal
	// TODO: only do this in debugging mode. Debug parameter can be set as input args
	debug.SendEventsToJournal(l.info.ContainerID, s.Text(), journal.PriInfo)
	// Send logs to aws cloudwatch logs
	err := l.LogWithRetry(s.Bytes(), time.Now())
	if err != nil {
		return errors.Wrap(s.Err(), "failed to send logs to cloudwatch")
	}

	return nil
}

// LogWithRetry sends logs to cloudwatch with retry.
func (l *logDriver) LogWithRetry(line []byte, logTimestamp time.Time) error {
	retryTimes := 0
	message := newMessage(line, l.info.ContainerID, logTimestamp)
	backoff := newBackoff()
	err := retry.RetryNWithBackoff(
		backoff,
		logger.LogRetryMaxAttempts,
		func() error {
			retryTimes += 1
			return l.stream.Log(message)
		})
	if err != nil {
		err = errors.Wrapf(err, "sending container logs to cloudwatch has been retried for %d times", retryTimes)
		return err
	}

	return nil
}

// newBackoff creates a new Backoff object.
func newBackoff() retry.Backoff {
	return retry.NewExponentialBackoff(
		logger.LogRetryMinBackoff,
		logger.LogRetryMaxBackoff,
		logger.LogRetryJitter,
		logger.LogRetryMultiple)
}

// newMessage creates a new logger message.
func newMessage(line []byte, source string, logTimestamp time.Time) *dockerlogger.Message {
	msg := dockerlogger.NewMessage()
	msg.Line = line
	msg.Source = source
	msg.Timestamp = logTimestamp

	return msg
}
