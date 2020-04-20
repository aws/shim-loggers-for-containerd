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

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"

	"github.com/containerd/containerd/runtime/v2/logging"
	"github.com/coreos/go-systemd/journal"
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
	stream, err := dockerawslogs.New(*info)
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
	err = l.Start(la.globalArgs.UID, la.globalArgs.GID, la.globalArgs.CleanupTime, ready)
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
