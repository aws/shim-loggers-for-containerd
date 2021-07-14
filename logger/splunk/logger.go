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

package splunk

import (
	"context"

	"github.com/containerd/containerd/runtime/v2/logging"
	dockersplunk "github.com/docker/docker/daemon/logger/splunk"
	"github.com/pkg/errors"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"
)

// splunk driver argument keys
const (
	// Required
	TokenKey = "splunk-token"
	URLKey   = "splunk-url"

	// Optional
	SourceKey             = "splunk-source"
	SourcetypeKey         = "splunk-sourcetype"
	IndexKey              = "splunk-index"
	CapathKey             = "splunk-capath"
	CanameKey             = "splunk-caname"
	InsecureskipverifyKey = "splunk-insecureskipverify"
	FormatKey             = "splunk-format"
	VerifyConnectionKey   = "splunk-verify-connection"
	GzipKey               = "splunk-gzip"
	GzipLevelKey          = "splunk-gzip-level"
	SplunkTagKey          = "splunk-tag"
	LabelsKey             = "labels"
	EnvKey                = "env"
	EnvRegexKey           = "env-regex"

	// Convert input parameter "splunk-tag" to the splunk parameter "tag"
	// This is to distinguish between the "tag" parameter from the fluentd input
	tagKey = "tag"
)

// Args represents splunk log driver arguments
type Args struct {
	Token              string
	URL                string
	Source             string
	Sourcetype         string
	Index              string
	Capath             string
	Caname             string
	Insecureskipverify string
	Format             string
	VerifyConnection   string
	Gzip               string
	GzipLevel          string
	Tag                string
	// TagSpecified represents whether a splunk tag was specified. It is used to differentiate between the default value
	// of the splunk tag when the flag is initialized and the client specifying the default value of the tag,
	// as they may be the same.
	TagSpecified bool
	Labels       string
	Env          string
	EnvRegex     string
}

// LoggerArgs stores global logger args and splunk specific args
type LoggerArgs struct {
	globalArgs    *logger.GlobalArgs
	dockerConfigs *logger.DockerConfigs
	args          *Args
}

// InitLogger initialize the input arguments
func InitLogger(globalArgs *logger.GlobalArgs, dockerConfigs *logger.DockerConfigs, splunkArgs *Args) *LoggerArgs {
	return &LoggerArgs{
		globalArgs:    globalArgs,
		dockerConfigs: dockerConfigs,
		args:          splunkArgs,
	}
}

// RunLogDriver initiates the splunk driver
func (la *LoggerArgs) RunLogDriver(ctx context.Context, config *logging.Config, ready func() error) error {
	defer debug.DeferFuncForRunLogDriver()

	loggerConfig := getSplunkConfig(la.args)
	info := logger.NewInfo(
		la.globalArgs.ContainerID,
		la.globalArgs.ContainerName,
		logger.WithConfig(loggerConfig),
	)
	info = logger.UpdateDockerConfigs(info, la.dockerConfigs)

	stream, err := dockersplunk.New(*info)
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
		debug.LoggerErr = errors.Wrap(err, "unable to create splunk log driver")
		return debug.LoggerErr
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToLog(logger.DaemonName, "Starting non-blocking mode driver", debug.INFO, 0)
		l = logger.NewBufferedLogger(l, la.globalArgs.MaxBufferSize, la.globalArgs.ContainerID)
	}

	// Start splunk log driver
	debug.SendEventsToLog(logger.DaemonName, "Starting splunk driver", debug.INFO, 0)
	err = l.Start(ctx, la.globalArgs.UID, la.globalArgs.GID, la.globalArgs.CleanupTime, ready)
	if err != nil {
		debug.LoggerErr = errors.Wrap(err, "failed to run splunk driver")
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

// getSplunkConfig sets values for splunk config
func getSplunkConfig(arg *Args) map[string]string {
	config := make(map[string]string)
	// Required arguments
	config[TokenKey] = arg.Token
	config[URLKey] = arg.URL
	// Optional arguments
	if arg.Source != "" {
		config[SourceKey] = arg.Source
	}
	if arg.Sourcetype != "" {
		config[SourcetypeKey] = arg.Sourcetype
	}
	if arg.Index != "" {
		config[IndexKey] = arg.Index
	}
	if arg.Capath != "" {
		config[CapathKey] = arg.Capath
	}
	if arg.Caname != "" {
		config[CanameKey] = arg.Caname
	}
	if arg.Insecureskipverify != "" {
		config[InsecureskipverifyKey] = arg.Insecureskipverify
	}
	if arg.Format != "" {
		config[FormatKey] = arg.Format
	}
	if arg.VerifyConnection != "" {
		config[VerifyConnectionKey] = arg.VerifyConnection
	}
	if arg.Gzip != "" {
		config[GzipKey] = arg.Gzip
	}
	if arg.GzipLevel != "" {
		config[GzipLevelKey] = arg.GzipLevel
	}
	if arg.TagSpecified {
		config[tagKey] = arg.Tag
	}
	if arg.Labels != "" {
		config[LabelsKey] = arg.Labels
	}
	if arg.Env != "" {
		config[EnvKey] = arg.Env
	}
	if arg.EnvRegex != "" {
		config[EnvRegexKey] = arg.EnvRegex
	}
	return config
}
