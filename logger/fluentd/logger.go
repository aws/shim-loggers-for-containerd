// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package fluentd

import (
	"context"
	"fmt"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"

	"github.com/containerd/containerd/runtime/v2/logging"
	dockerfluentd "github.com/docker/docker/daemon/logger/fluentd"
)

const (
	AddressKey            = "fluentd-address"
	AsyncConnectKey       = "fluentd-async-connect"
	FluentdTagKey         = "fluentd-tag"
	SubsecondPrecisionKey = "fluentd-sub-second-precision"
	BufferLimitKey        = "fluentd-buffer-limit"

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
	BufferLimit        string
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
		debug.LoggerErr = fmt.Errorf("unable to create stream: %w", err)
		return debug.LoggerErr
	}

	l, err := logger.NewLogger(
		logger.WithStdout(config.Stdout),
		logger.WithStderr(config.Stderr),
		logger.WithInfo(info),
		logger.WithStream(stream),
	)
	if err != nil {
		debug.LoggerErr = fmt.Errorf("unable to create fluentd driver: %w", err)
		return debug.LoggerErr
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToLog(logger.DaemonName, "Starting non-blocking mode driver", debug.INFO, 0)
		l = logger.NewBufferedLogger(l, logger.DefaultBufSizeInBytes, la.globalArgs.MaxBufferSize, la.globalArgs.ContainerID)
	}

	// Start fluentd driver
	debug.SendEventsToLog(logger.DaemonName, "Starting fluentd driver", debug.INFO, 0)
	err = l.Start(ctx, la.globalArgs.UID, la.globalArgs.GID, la.globalArgs.CleanupTime, ready)
	if err != nil {
		debug.LoggerErr = fmt.Errorf("failed to run fluentd driver: %w", err)
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
	config[BufferLimitKey] = args.BufferLimit

	return config
}
