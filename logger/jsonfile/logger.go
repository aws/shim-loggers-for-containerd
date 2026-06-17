// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package jsonfile provides functionalities for integrating the json-file logging driver
// with shim-loggers-for-containerd.
package jsonfile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/runtime/v2/logging"
	dockerlogger "github.com/docker/docker/daemon/logger"
	dockerjsonfilelog "github.com/docker/docker/daemon/logger/jsonfilelog"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"
)

// json-file driver argument keys.
const (
	// DriverName is the name of the json-file log driver.
	DriverName = "json-file"

	// LogPathKey specifies the per-container output file path. The shim-logger creates
	// the parent directory if it does not exist.
	LogPathKey = "log-path"

	// MaxSizeKey is the maximum size of the log file before it is rolled (e.g., "10m").
	// Forwarded to moby's jsonfilelog as-is. moby rejects "<= 0" or malformed values when
	// the writer starts up.
	MaxSizeKey = "max-size"

	// MaxFileKey is the maximum number of log files that can be present (e.g., "5").
	// Forwarded to moby's jsonfilelog as-is.
	MaxFileKey = "max-file"

	// CompressKey indicates whether to compress rotated log files (gzip). Defaults to false
	// in moby. Note: moby rejects compress=true when max-file < 2 or max-size is unset.
	CompressKey = "compress"

	// LabelsKey is the moby-side option key for label-based extras (renamed from
	// JSONFileLabelsKey on the way into moby). The shim-logger's input flag is prefixed
	// to avoid collisions with the splunk driver's "labels" flag.
	LabelsKey = "labels"

	// LabelsRegexKey is the moby-side option key for label-regex extras (renamed from
	// JSONFileLabelsRegexKey on the way into moby).
	LabelsRegexKey = "labels-regex"

	// EnvKey is the moby-side option key for env-based extras (renamed from JSONFileEnvKey).
	EnvKey = "env"

	// EnvRegexKey is the moby-side option key for env-regex extras (renamed from
	// JSONFileEnvRegexKey).
	EnvRegexKey = "env-regex"

	// JSONFileLabelsKey is the input parameter name for the labels list. Renamed on
	// the way into moby to "labels" so the input parameter doesn't collide with the
	// same-named parameter from the splunk driver.
	JSONFileLabelsKey = "json-file-labels"

	// JSONFileLabelsRegexKey is the input parameter name for the labels regex.
	JSONFileLabelsRegexKey = "json-file-labels-regex"

	// JSONFileEnvKey is the input parameter name for the env list.
	JSONFileEnvKey = "json-file-env"

	// JSONFileEnvRegexKey is the input parameter name for the env regex.
	JSONFileEnvRegexKey = "json-file-env-regex"

	// JSONFileTagKey is the input parameter name for the json-file tag template.
	// Renamed on the way into moby to "tag" so the input parameter doesn't collide with
	// the same-named parameters from the splunk and fluentd drivers.
	JSONFileTagKey = "json-file-tag"

	// tagKey is the moby-side option key for tag template (renamed from JSONFileTagKey).
	tagKey = "tag"
)

// logDirMode is the permission mode for per-container log directories.
// Setgid (2000) ensures new files inherit the parent directory's group.
// Owner=rwx, group=r-x, other=none.
const logDirMode = os.FileMode(02750)

// Args represents json-file log driver arguments.
type Args struct {
	// Required.
	LogPath string

	// Optional.
	MaxSize     string
	MaxFile     string
	Compress    string
	Labels      string
	LabelsRegex string
	Env         string
	EnvRegex    string
	Tag         string
	// TagSpecified represents whether a json-file tag was specified. Used to differentiate
	// between the default empty string and an explicitly-set empty tag, mirroring the splunk
	// driver's pattern.
	TagSpecified bool
}

// LoggerArgs stores global logger args and json-file specific args.
type LoggerArgs struct {
	globalArgs *logger.GlobalArgs
	args       *Args
}

// InitLogger initializes the input arguments.
func InitLogger(globalArgs *logger.GlobalArgs, jsonFileArgs *Args) *LoggerArgs {
	return &LoggerArgs{
		globalArgs: globalArgs,
		args:       jsonFileArgs,
	}
}

// RunLogDriver initializes and starts the json-file logger.
// Errors with the log driver are logged but not returned to prevent container termination.
func (la *LoggerArgs) RunLogDriver(ctx context.Context, config *logging.Config, ready func() error) error {
	defer debug.DeferFuncForRunLogDriver()

	loggerConfig, err := getJSONFileConfig(la.args)
	if err != nil {
		debug.ErrLogger = fmt.Errorf("unable to validate log options: %w", err)
		return debug.ErrLogger
	}
	info := logger.NewInfo(
		la.globalArgs.ContainerID,
		la.globalArgs.ContainerName,
		logger.WithConfig(loggerConfig),
		logger.WithLogPath(la.args.LogPath),
	)

	// Create the log file's parent directory if it does not exist.
	if dir := filepath.Dir(la.args.LogPath); dir != "" {
		if err := os.MkdirAll(dir, logDirMode); err != nil {
			debug.ErrLogger = fmt.Errorf("unable to create log directory %s: %w", dir, err)
			return debug.ErrLogger
		}
	}

	stream, err := dockerjsonfilelog.New(*info)
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
		debug.ErrLogger = fmt.Errorf("unable to create json-file driver: %w", err)
		return debug.ErrLogger
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToLog(logger.DaemonName, "Starting non-blocking mode driver", debug.INFO, 0)
		l = logger.NewBufferedLogger(l, logger.DefaultBufSizeInBytes, la.globalArgs.MaxBufferSize, la.globalArgs.ContainerID)
	}

	// Start json-file driver.
	debug.SendEventsToLog(logger.DaemonName, "Starting json-file driver", debug.INFO, 0)
	err = l.Start(ctx, la.globalArgs.CleanupTime, ready)
	if err != nil {
		debug.ErrLogger = fmt.Errorf("failed to run json-file driver: %w", err)
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

// getJSONFileConfig sets values for json-file config and validates them via moby's
// upstream ValidateLogOpts. Optional fields are only set when non-empty so we don't
// trigger moby's "unknown log opt" rejection on empty strings.
func getJSONFileConfig(arg *Args) (map[string]string, error) {
	config := make(map[string]string)
	if arg.MaxSize != "" {
		config[MaxSizeKey] = arg.MaxSize
	}
	if arg.MaxFile != "" {
		config[MaxFileKey] = arg.MaxFile
	}
	if arg.Compress != "" {
		config[CompressKey] = arg.Compress
	}
	if arg.Labels != "" {
		config[LabelsKey] = arg.Labels
	}
	if arg.LabelsRegex != "" {
		config[LabelsRegexKey] = arg.LabelsRegex
	}
	if arg.Env != "" {
		config[EnvKey] = arg.Env
	}
	if arg.EnvRegex != "" {
		config[EnvRegexKey] = arg.EnvRegex
	}
	if arg.TagSpecified {
		config[tagKey] = arg.Tag
	}

	err := dockerlogger.ValidateLogOpts(DriverName, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
