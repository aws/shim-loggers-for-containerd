// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package awslogs provides functionalities for integrating the awslogs logging driver
// with shim-loggers-for-containerd.
package awslogs

import (
	"context"
	"fmt"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"

	"github.com/containerd/containerd/runtime/v2/logging"
	dockerawslogs "github.com/docker/docker/daemon/logger/awslogs"
)

const (
	// awslogs driver options.

	// RegionKey specifies the AWS logging region.
	RegionKey = "awslogs-region"
	// GroupKey denotes the AWS logging group name.
	GroupKey = "awslogs-group"
	// CreateGroupKey is a flag to create a new AWS logging group.
	CreateGroupKey = "awslogs-create-group"
	// StreamKey specifies the AWS logging stream name.
	StreamKey = "awslogs-stream"
	// CreateStreamKey is a flag to create a new log stream.
	CreateStreamKey = "awslogs-create-stream"
	// MultilinePatternKey defines the pattern for multiline logs.
	MultilinePatternKey = "awslogs-multiline-pattern"
	// DatetimeFormatKey specifies the datetime format of the logs.
	DatetimeFormatKey = "awslogs-datetime-format"
	// CredentialsEndpointKey denotes the AWS credentials endpoint (not actual credentials).
	CredentialsEndpointKey = "awslogs-credentials-endpoint" //nolint:gosec // not credentials
	// EndpointKey is the AWS logging endpoint.
	EndpointKey = "awslogs-endpoint"
	// LogFormatKey is used to explicitly set EMF header.
	LogFormatKey = "awslogs-format"
	// JSONEmfLogFormat currently only 'json/emf' is supported.
	// See: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch_Embedded_Metric_Format_Specification.html
	JSONEmfLogFormat = "json/emf"

	// There are 26 bytes additional bytes for each log event:
	// See more details in: http://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_PutLogEvents.html
	perEventBytes = 26
	// The value of maximumBytesPerEvent is adopted from Docker. Reference:
	// https://github.com/moby/moby/blob/19.03/daemon/logger/awslogs/cloudwatchlogs.go#L58
	maximumBytesPerEvent = 262144 - perEventBytes

	// The max size of CloudWatch events is 256kb.
	defaultAwsBufSizeInBytes = 256 * 1024
)

// Args represents AWSlogs driver arguments.
type Args struct {
	// Required arguments.
	Group               string
	Region              string
	Stream              string
	CredentialsEndpoint string

	// Optional arguments.
	CreateGroup      string
	CreateStream     string
	MultilinePattern string
	DatetimeFormat   string
	Endpoint         string
	LogsFormatHeader string
}

// LoggerArgs stores global logger args and awslogs specific args.
type LoggerArgs struct {
	globalArgs *logger.GlobalArgs
	args       *Args
}

// InitLogger initialize the input arguments.
func InitLogger(globalArgs *logger.GlobalArgs, awslogsArgs *Args) *LoggerArgs {
	return &LoggerArgs{
		globalArgs: globalArgs,
		args:       awslogsArgs,
	}
}

// RunLogDriver initiates an awslogs driver and starts driving container logs to cloudwatch.
func (la *LoggerArgs) RunLogDriver(ctx context.Context, config *logging.Config, ready func() error) error {
	defer debug.DeferFuncForRunLogDriver()

	loggerConfig := getAWSLogsConfig(la.args)
	if err := validateLogOptCompatability(loggerConfig); err != nil {
		debug.ErrLogger = fmt.Errorf("incompatible logger options: %w", err)
		return debug.ErrLogger
	}
	info := logger.NewInfo(
		la.globalArgs.ContainerID,
		la.globalArgs.ContainerName,
		logger.WithConfig(loggerConfig),
	)
	stream, err := dockerawslogs.New(*info)
	if err != nil {
		debug.ErrLogger = fmt.Errorf("unable to create stream: %w", err)
		return debug.ErrLogger
	}

	l, err := logger.NewLogger(
		logger.WithStdout(config.Stdout),
		logger.WithStderr(config.Stderr),
		logger.WithInfo(info),
		logger.WithStream(stream),
		logger.WithBufferSizeInBytes(maximumBytesPerEvent),
	)
	if err != nil {
		debug.ErrLogger = fmt.Errorf("unable to create awslogs driver: %w", err)
		return debug.ErrLogger
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToLog(logger.DaemonName, "Starting log streaming for non-blocking mode awslogs driver",
			debug.INFO, 0)
		l = logger.NewBufferedLogger(l, defaultAwsBufSizeInBytes, la.globalArgs.MaxBufferSize, la.globalArgs.ContainerID)
	}

	// Start awslogs driver
	debug.SendEventsToLog(logger.DaemonName, "Starting log streaming for awslogs driver", debug.INFO, 0)
	err = l.Start(ctx, la.globalArgs.CleanupTime, ready)
	if err != nil {
		debug.ErrLogger = fmt.Errorf("failed to run awslogs driver: %w", err)
		// Do not return error if log driver has issue sending logs to destination, because if error
		// returned here, containerd will identify this error and kill shim process, which will kill
		// the container process accordingly.
		// Note: the container will continue to run if shim logger exits after.
		// Reference: https://github.com/containerd/containerd/blob/release/1.3/runtime/v2/logging/logging.go
		return nil
	}
	debug.SendEventsToLog(logger.DaemonName, "Logging finished", debug.INFO, 1)

	return nil
}

// getAWSLogsConfig sets values for awslogs config.
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
	createStream := args.CreateStream
	if createStream != "" {
		config[CreateStreamKey] = createStream
	}
	multilinePattern := args.MultilinePattern
	if multilinePattern != "" {
		config[MultilinePatternKey] = multilinePattern
	}
	datetimeFormat := args.DatetimeFormat
	if datetimeFormat != "" {
		config[DatetimeFormatKey] = datetimeFormat
	}
	endpoint := args.Endpoint
	if endpoint != "" {
		config[EndpointKey] = endpoint
	}
	logsFormatHeader := args.LogsFormatHeader
	if logsFormatHeader != "" {
		config[LogFormatKey] = logsFormatHeader
	}

	return config
}

func validateLogOptCompatability(cfg map[string]string) error {
	_, datetimeFormatKeyExists := cfg[DatetimeFormatKey]
	_, multilinePatternKeyExists := cfg[MultilinePatternKey]

	if cfg[LogFormatKey] != "" {
		// Only json/emf is supported at the moment.
		if cfg[LogFormatKey] != JSONEmfLogFormat {
			return fmt.Errorf("unsupported log format '%s'", cfg[LogFormatKey])
		}
		if datetimeFormatKeyExists || multilinePatternKeyExists {
			return fmt.Errorf(
				"you cannot configure log opt '%s' or '%s' when log opt '%s' is set to '%s'",
				DatetimeFormatKey,
				MultilinePatternKey,
				LogFormatKey,
				JSONEmfLogFormat,
			)
		}
	}

	return nil
}
