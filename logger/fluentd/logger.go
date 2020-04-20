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

package fluentd

import (
	"context"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"

	"github.com/containerd/containerd/runtime/v2/logging"
	"github.com/coreos/go-systemd/journal"
	dockerfluentd "github.com/docker/docker/daemon/logger/fluentd"
	"github.com/pkg/errors"
)

const (
	AddressKey      = "fluentd-address"
	AsyncConnectKey = "fluentd-async-connect"
	FluentdTagKey   = "fluentd-tag"

	// Convert input parameter "fluentd-tag" to the fluentd parameter "tag"
	// This is to distinguish between the "tag" parameter from the splunk input
	tagKey = "tag"
)

// Args represents fluentd log driver arguments
type Args struct {
	// Optional arguments
	Address      string
	AsyncConnect string
	Tag          string
}

// LoggerArgs stores global logger args and fluentd specific args
type LoggerArgs struct {
	globalArgs *logger.GlobalArgs
	args       *Args
}

// InitLogger initialize the input arguments
func InitLogger(globalArgs *logger.GlobalArgs, fluentdArgs *Args) *LoggerArgs {
	return &LoggerArgs{
		globalArgs: globalArgs,
		args:       fluentdArgs,
	}
}

func (la *LoggerArgs) RunLogDriver(ctx context.Context, config *logging.Config, ready func() error) error {
	loggerConfig := getFluentdConfig(la.args)
	info := logger.NewInfo(
		la.globalArgs.ContainerID,
		la.globalArgs.ContainerName,
		logger.WithConfig(loggerConfig),
	)
	stream, err := dockerfluentd.New(*info)
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
		return errors.Wrap(err, "unable to create fluentd driver")
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToJournal(logger.DaemonName, "Starting non-blocking mode driver", journal.PriInfo)
		l = logger.NewBufferedLogger(l, la.globalArgs.MaxBufferSize)
	}

	// Start fluentd driver
	debug.SendEventsToJournal(logger.DaemonName, "Starting fluentd driver", journal.PriInfo)
	err = l.Start(la.globalArgs.UID, la.globalArgs.GID, la.globalArgs.CleanupTime, ready)
	if err != nil {
		return errors.Wrap(err, "failed to start fluentd driver")
	}
	debug.SendEventsToJournal(logger.DaemonName, "Logging finished", journal.PriInfo)

	return nil
}

// getFluentdConfig sets values for fluentd config
func getFluentdConfig(args *Args) map[string]string {
	config := make(map[string]string)
	config[tagKey] = args.Tag
	config[AddressKey] = args.Address
	config[AsyncConnectKey] = args.AsyncConnect

	return config
}
