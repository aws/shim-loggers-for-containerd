package awslogs

import (
	"context"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"

	apierrors "github.com/aws/amazon-ecs-agent/agent/api/errors"
	"github.com/aws/amazon-ecs-agent/agent/utils/retry"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/containerd/containerd/runtime/v2/logging"
	"github.com/coreos/go-systemd/journal"
	dockerlogger "github.com/docker/docker/daemon/logger"
	dockerawslogs "github.com/docker/docker/daemon/logger/awslogs"
	"github.com/pkg/errors"
)

const (
	// awslogs driver options
	RegionKey              = "awslogs-region"
	GroupKey               = "awslogs-group"
	CreateGroupKey         = "awslogs-create-group"
	StreamKey              = "awslogs-stream"
	MultilinePatternKey    = "awslogs-multiline-pattern"
	DatetimeFormatKey      = "awslogs-datetime-format"
	CredentialsEndpointKey = "awslogs-credentials-endpoint"

	// Define the retry parameters for retrying creating stream
	createStreamRetryMaxAttempts = 5
	createStreamRetryMinBackoff  = 500 * time.Millisecond
	createStreamRetryMaxBackoff  = 2 * time.Second
	createStreamRetryJitter      = 1
	createStreamRetryMultiple    = 2
)

// Args represents AWSlogs driver arguments
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

// loggerStream is a type of function of docker new logger API for awslogs
// log driver. This is defined mostly for testing purpose.
type loggerStream func(dockerlogger.Info) (dockerlogger.Logger, error)

// InitLogger initialize the input arguments
func InitLogger(globalArgs *logger.GlobalArgs, awslogsArgs *Args) *LoggerArgs {
	return &LoggerArgs{
		globalArgs: globalArgs,
		args:       awslogsArgs,
	}
}

// RunLogDriver initiates an awslogs driver and starts driving container logs to cloudwatch
func (la *LoggerArgs) RunLogDriver(ctx context.Context, config *logging.Config, ready func() error) error {
	loggerConfig := getAWSLogsConfig(la.args)
	info := logger.NewInfo(
		la.globalArgs.ContainerID,
		la.globalArgs.ContainerName,
		logger.WithConfig(loggerConfig),
	)
	stream, err := createStreamWithRetry(dockerawslogs.New, info)
	if err != nil {
		return errors.Wrap(err, "unable to create stream")
	}

	l, err := logger.NewLogger(
		logger.WithStdout(config.Stdout),
		logger.WithStderr(config.Stderr),
		logger.WithInfo(info),
		logger.WithStream(stream),
	)
	if err != nil {
		return errors.Wrap(err, "unable to create awslogs driver")
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToJournal(logger.DaemonName, "Starting non-blocking mode driver", journal.PriInfo)
		l = logger.NewBufferedLogger(l, la.globalArgs.MaxBufferSize)
	}

	// Start awslogs driver
	debug.SendEventsToJournal(logger.DaemonName, "Starting awslogs driver", journal.PriInfo)
	err = l.Start(la.globalArgs.UID, la.globalArgs.GID, ready)
	if err != nil {
		return errors.Wrap(err, "failed to start awslogs driver")
	}
	debug.SendEventsToJournal(logger.DaemonName, "Logging finished", journal.PriInfo)

	return nil
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

// createStreamWithRetry creates log stream with reties if the returned error is retriable error.
func createStreamWithRetry(createStream loggerStream, info *dockerlogger.Info) (logger.Client, error) {
	var (
		stream     logger.Client
		cErr       error
		retryTimes int
	)
	backoff := newCreateStreamBackoff()
	err := retry.RetryNWithBackoff(
		backoff,
		createStreamRetryMaxAttempts,
		func() error {
			retryTimes += 1
			stream, cErr = createStreamOnce(createStream, info)
			return cErr
		})
	if err != nil {
		return nil, errors.Wrapf(err, "creating log stream has been retried %d time(s)", retryTimes)
	}

	return stream, nil
}

// createStreamOnce creates log stream and return an retriable error.
func createStreamOnce(createStream loggerStream, info *dockerlogger.Info) (logger.Client, error) {
	var (
		stream logger.Client
		err    error
	)
	stream, err = createStream(*info)
	if err == nil {
		return stream, nil
	}

	return nil, createStreamRetriableError(errors.Cause(err))
}

// createStreamRetriableError checks if error returned by cloudwatch logs CreateLogGroup or CreateLogStream
// API is retriable and returns corresponding error. For now we only retry on error of OperationAbortedException
// and ResourceAlreadyExistsException. Note that ResourceAlreadyExistsException is already checked by docker:
//		https://github.com/moby/moby/blob/master/daemon/logger/awslogs/cloudwatchlogs.go#L468
// We do double check here for safety.
// Reference: https://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_CreateLogGroup.html
func createStreamRetriableError(err error) error {
	retriable := false
	aerr, ok := err.(awserr.Error)
	if ok {
		switch aerr.Code() {
		case cloudwatchlogs.ErrCodeOperationAbortedException:
			retriable = true
		case cloudwatchlogs.ErrCodeResourceAlreadyExistsException:
			retriable = true
		}
	}

	return apierrors.NewRetriableError(apierrors.NewRetriable(retriable), err)
}

// newCreateStreamBackoff creates a new Backoff object used for creating stream.
func newCreateStreamBackoff() retry.Backoff {
	return retry.NewExponentialBackoff(
		createStreamRetryMinBackoff,
		createStreamRetryMaxBackoff,
		createStreamRetryJitter,
		createStreamRetryMultiple)
}
