// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package fluentd provides functionalities for integrating the fluentd logging driver
// with shim-loggers-for-containerd.
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
	// AddressKey specifies the address configuration for fluentd.
	AddressKey = "fluentd-address"
	// AsyncConnectKey specifies the async connect configuration key for fluentd.
	AsyncConnectKey = "fluentd-async"
	// FluentdTagKey specifies the tag configuration key for fluentd.
	FluentdTagKey = "fluentd-tag"
	// SubsecondPrecisionKey specifies the sub-second precision configuration key for fluentd.
	SubsecondPrecisionKey = "fluentd-sub-second-precision"
	// BufferLimitKey specifies the buffer limit configuration key for fluentd.
	BufferLimitKey = "fluentd-buffer-limit"
	// WriteTimeoutKey specifies the write timeout configuration key for fluentd.
	WriteTimeoutKey = "fluentd-write-timeout"

	// Convert input parameter "fluentd-tag" to the fluentd parameter "tag".
	// This is to distinguish between the "tag" parameter from the splunk input.
	tagKey = "tag"
)

// Args represents fluentd log driver arguments.
type Args struct {
	// Optional arguments
	Address            string
	AsyncConnect       string
	Tag                string
	SubsecondPrecision string
	BufferLimit        string
	WriteTimeout       string
}

// LoggerArgs stores global logger args and fluentd specific args.
type LoggerArgs struct {
	globalArgs *logger.GlobalArgs
	args       *Args
}

// InitLogger initialize the input arguments.
func InitLogger(globalArgs *logger.GlobalArgs, fluentdArgs *Args) *LoggerArgs {
	return &LoggerArgs{
		globalArgs: globalArgs,
		args:       fluentdArgs,
	}
}

// RunLogDriver initializes and starts the fluentd logger.
// Errors with the log driver are logged but not returned to prevent container termination.
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
		debug.ErrLogger = fmt.Errorf("unable to create stream: %w", err)
		return debug.ErrLogger
	}

	l, err := logger.NewLogger(
		logger.WithStdout(config.Stdout),
		logger.WithStderr(config.Stderr),
		logger.WithInfo(info),
		logger.WithStream(stream),
	)
	if err != nil {
		debug.ErrLogger = fmt.Errorf("unable to create fluentd driver: %w", err)
		return debug.ErrLogger
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToLog(logger.DaemonName, "Starting non-blocking mode driver", debug.INFO, 0)
		l = logger.NewBufferedLogger(l, logger.DefaultBufSizeInBytes, la.globalArgs.MaxBufferSize, la.globalArgs.ContainerID)
	}

	// Start fluentd driver
	debug.SendEventsToLog(logger.DaemonName, "Starting fluentd driver", debug.INFO, 0)
	err = l.Start(ctx, la.globalArgs.CleanupTime, ready)
	if err != nil {
		debug.ErrLogger = fmt.Errorf("failed to run fluentd driver: %w", err)
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

// getFluentdConfig sets values for fluentd config.
func getFluentdConfig(args *Args) map[string]string {
	config := make(map[string]string)
	config[tagKey] = args.Tag
	config[AddressKey] = args.Address
	config[AsyncConnectKey] = args.AsyncConnect
	config[SubsecondPrecisionKey] = args.SubsecondPrecision
	config[BufferLimitKey] = args.BufferLimit
	config[WriteTimeoutKey] = args.WriteTimeout

	return config
}
