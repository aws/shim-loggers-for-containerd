package main

import (
	"github.com/aws/shim-loggers-for-containerd/logger/awslogs"
	"github.com/aws/shim-loggers-for-containerd/logger/fluentd"

	"github.com/spf13/pflag"
)

const (
	// Container options
	containerIDKey   = "container-id"
	containerNameKey = "container-name"

	// Mode and buffer size options
	modeKey          = "mode"
	maxBufferSizeKey = "max-buffer-size"

	// LogDriver options
	logDriverTypeKey  = "log-driver"
	awslogsDriverName = "awslogs"
	fluentdDriverName = "fluentd"

	// Verbose mode option
	verboseKey = "verbose"
)

// initCommonLogOpts initialize common options that get used by any log drivers
func initCommonLogOpts() {
	// container info
	pflag.String(containerIDKey, "", "Id of the container")
	pflag.String(containerNameKey, "", "Name of the container")

	// log driver options
	pflag.String(logDriverTypeKey, "", "`awslogs`, `fluentd` or `splunk`")

	// mode options
	pflag.String(modeKey, "", "Whether the writer is blocked or not blocked")
	pflag.String(maxBufferSizeKey, "", "The size of intermediate buffer for non-blocking mode")

	// verbose mode option
	pflag.Bool(verboseKey, false, "If set, then more logs will be printed for debugging")
}

// initAWSLogsOpts initialize awslogs driver specified options
func initAWSLogsOpts() {
	pflag.String(awslogs.GroupKey, "", "The CloudWatch log group to use")
	pflag.String(awslogs.RegionKey, "", "The CloudWatch region to use")
	pflag.String(awslogs.StreamKey, "", "The CloudWatch log stream to use")
	pflag.String(awslogs.CreateGroupKey, "", "Is this a new group that needs to be created?")
	pflag.String(awslogs.CredentialsEndpointKey, "", "The endpoint for iam credentials")
	pflag.String(awslogs.MultilinePatternKey, "", "Support multiline pattern for debug")
	pflag.String(awslogs.DatetimeFormatKey, "", "Multiline pattern in strftime format")
}

// initFluentdOpts initialize fluentd driver specified options
func initFluentdOpts() {
	pflag.String(fluentd.AddressKey, "", "The address connected to Fluentd daemon")
	pflag.Bool(fluentd.AsyncConnectKey, false, "If connecting Fluentd daemon in background")
	pflag.String(fluentd.TagKey, "", "The tag used to identify log messages")
}
