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

package awslogs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"

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
	CreateStreamKey        = "awslogs-create-stream"
	StreamKey              = "awslogs-stream"
	MultilinePatternKey    = "awslogs-multiline-pattern"
	DatetimeFormatKey      = "awslogs-datetime-format"
	CredentialsEndpointKey = "awslogs-credentials-endpoint"

	// There are 26 bytes additional bytes for each log event:
	// See more details in: http://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_PutLogEvents.html
	perEventBytes = 26
	// The value of maximumBytesPerEvent is adopted from Docker. Reference:
	// https://github.com/moby/moby/blob/19.03/daemon/logger/awslogs/cloudwatchlogs.go#L58
	maximumBytesPerEvent = 262144 - perEventBytes
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
	CreateStream     string
	MultilinePattern string
	DatetimeFormat   string
}

// LoggerArgs stores global logger args and awslogs specific args
type LoggerArgs struct {
	globalArgs *logger.GlobalArgs
	args       *Args
}

// InitLogger initialize the input arguments
func InitLogger(globalArgs *logger.GlobalArgs, awslogsArgs *Args) *LoggerArgs {
	return &LoggerArgs{
		globalArgs: globalArgs,
		args:       awslogsArgs,
	}
}

// RunLogDriver initiates an awslogs driver and starts driving container logs to cloudwatch
func (la *LoggerArgs) RunLogDriver(ctx context.Context, config *logging.Config, ready func() error) error {
	defer debug.DeferFuncForRunLogDriver()

	loggerConfig := getAWSLogsConfig(la.args)
	info := logger.NewInfo(
		la.globalArgs.ContainerID,
		la.globalArgs.ContainerName,
		logger.WithConfig(loggerConfig),
	)

	stream, err := la.initialize(info)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "failed to initialize awslogs driver")
		return debug.LoggerErr
	}

	l, err := logger.NewLogger(
		logger.WithStdout(config.Stdout),
		logger.WithStderr(config.Stderr),
		logger.WithInfo(info),
		logger.WithStream(stream),
		logger.WithBufferSizeInBytes(maximumBytesPerEvent),
	)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "unable to create awslogs driver")
		return debug.LoggerErr
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToJournal(logger.DaemonName, "Starting log streaming for non-blocking mode awslogs driver",
			journal.PriInfo, 0)
		l = logger.NewBufferedLogger(l, la.globalArgs.MaxBufferSize, la.globalArgs.ContainerID)
	}

	// Start awslogs driver
	debug.SendEventsToJournal(logger.DaemonName, "Starting log streaming for awslogs driver", journal.PriInfo, 0)
	err = l.Start(ctx, la.globalArgs.UID, la.globalArgs.GID, la.globalArgs.CleanupTime, ready)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "failed to run awslogs driver")
		// Do not return error if log driver has issue sending logs to destination, because if error
		// returned here, containerd will identify this error and kill shim process, which will kill
		// the container process accordingly.
		// Note: the container will continue to run if shim logger exits after.
		// Reference: https://github.com/containerd/containerd/blob/release/1.3/runtime/v2/logging/logging.go
		return nil
	}
	debug.SendEventsToJournal(logger.DaemonName, "Logging finished", journal.PriInfo, 1)

	return nil
}

// initialize creates log stream config, CloudWatch Logs client, and log stream in CloudWatch Logs
// if specified. It also starts a routine to collect batch before container starts streaming logs.
// It breaks down the steps in awslogs logger New function from moby:
// https://github.com/moby/moby/blob/ad1b781e44fa1e44b9e654e5078929aec56aed66/daemon/logger/awslogs/cloudwatchlogs.go#L140
// This is added as a work-around to make CreateLogStream API call optional.
func (la *LoggerArgs) initialize(info *dockerlogger.Info) (*dockerawslogs.LogStream, error) {
	containerStreamConfig, err := dockerawslogs.NewStreamConfig(*info)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "unable to create CloudWatch stream config")
		return nil, debug.LoggerErr
	}

	debug.SendEventsToJournal(logger.DaemonName, "Creating CloudWatch logs client", journal.PriInfo, 0)
	client, err := dockerawslogs.NewAWSLogsClient(*info)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "unable to create CloudWatch client for awslogs driver")
		return nil, debug.LoggerErr
	}
	debug.SendEventsToJournal(logger.DaemonName, "CloudWatch logs client created", journal.PriInfo, 0)

	containerStream := &dockerawslogs.LogStream{
		LogStreamName:    containerStreamConfig.LogStreamName,
		LogGroupName:     containerStreamConfig.LogGroupName,
		LogCreateGroup:   containerStreamConfig.LogCreateGroup,
		MultilinePattern: containerStreamConfig.MultilinePattern,
		Client:           client,
		Messages:         make(chan *dockerlogger.Message, dockerawslogs.DefaultMaxBufferedEvents),
	}

	// This channel is a required input parameter for CollectBatch function.
	creationDone := make(chan bool)
	err = la.createLogStream(containerStream)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "unable to create CloudWatch log stream")
		return nil, debug.LoggerErr
	}
	close(creationDone)

	// Start collecting batches and send log events by evoking PutLogEvents API calls.
	go containerStream.CollectBatch(creationDone)

	return containerStream, nil
}

// createLogStream creates log stream in CloudWatch Logs for container if required.
func (la *LoggerArgs) createLogStream(containerStream *dockerawslogs.LogStream) error {
	createStream, err := strconv.ParseBool(la.args.CreateStream)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "unable to parse create log stream boolean value")
		return debug.LoggerErr
	}
	if !createStream {
		debug.SendEventsToJournal(logger.DaemonName, "Skipping log stream creation", journal.PriInfo, 0)
		return nil
	}

	debug.SendEventsToJournal(logger.DaemonName,
		fmt.Sprintf("Creating log group %s and log stream %s",
			containerStream.LogGroupName,
			containerStream.LogStreamName),
		journal.PriInfo, 0)
	err = containerStream.Create()
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "unable to create stream in CloudWatch")
		return debug.LoggerErr
	}
	debug.SendEventsToJournal(logger.DaemonName, "Log group and log stream created", journal.PriInfo, 0)

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
