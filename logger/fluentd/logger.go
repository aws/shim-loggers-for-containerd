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
	dockerfluentd "github.com/docker/docker/daemon/logger/fluentd"
	"github.com/pkg/errors"
)

const (
	AddressKey            = "fluentd-address"
	AsyncConnectKey       = "fluentd-async-connect"
	FluentdTagKey         = "fluentd-tag"
	SubsecondPrecisionKey = "fluentd-sub-second-precision"

	// Convert input parameter "fluentd-tag" to the fluentd parameter "tag"
	// This is to distinguish between the "tag" parameter from the splunk input
	tagKey = "tag"
)

// Args represents fluentd log driver arguments
type Args struct {
	// Optional arguments
	Address            string
	AsyncConnect       string
	Tag                string
	SubsecondPrecision string
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
	defer debug.DeferFuncForRunLogDriver()

	loggerConfig := getFluentdConfig(la.args)
	info := logger.NewInfo(
		la.globalArgs.ContainerID,
		la.globalArgs.ContainerName,
		logger.WithConfig(loggerConfig),
	)
	stream, err := dockerfluentd.New(*info)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "unable to create stream")
		return debug.LoggerErr
	}

	l, err := logger.NewLogger(
		logger.WithStdout(config.Stdout),
		logger.WithStderr(config.Stderr),
		logger.WithInfo(info),
		logger.WithStream(stream),
	)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "unable to create fluentd driver")
		return debug.LoggerErr
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToLog(logger.DaemonName, "Starting non-blocking mode driver", debug.INFO, 0)
		l = logger.NewBufferedLogger(l, la.globalArgs.MaxBufferSize, la.globalArgs.ContainerID)
	}

	// Start fluentd driver
	debug.SendEventsToLog(logger.DaemonName, "Starting fluentd driver", debug.INFO, 0)
	err = l.Start(ctx, la.globalArgs.UID, la.globalArgs.GID, la.globalArgs.CleanupTime, ready)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "failed to run fluentd driver")
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

// getFluentdConfig sets values for fluentd config
func getFluentdConfig(args *Args) map[string]string {
	config := make(map[string]string)
	config[tagKey] = args.Tag
	config[AddressKey] = args.Address
	config[AsyncConnectKey] = args.AsyncConnect
	config[SubsecondPrecisionKey] = args.SubsecondPrecision

	return config
}
