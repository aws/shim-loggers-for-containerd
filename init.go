package main

import (
	"github.com/aws/shim-loggers-for-containerd/logger/awslogs"
	"github.com/aws/shim-loggers-for-containerd/logger/fluentd"
	"github.com/aws/shim-loggers-for-containerd/logger/splunk"

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
	splunkDriverName  = "splunk"

	// Verbose mode option
	verboseKey = "verbose"

	// UID/GID option
	uidKey = "uid"
	gidKey = "gid"
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

	// set uid/gid option
	pflag.Int(uidKey, -1, "Customized uid for all the goroutines in shim logger process")
	pflag.Int(gidKey, -1, "Customized gid for all the goroutines in shim logger process")
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
	pflag.String(fluentd.FluentdTagKey, "", "The tag used to identify log messages")
}

// initSplunkOpts initialize splunk driver specified options
// Argument usage taken from https://docs.docker.com/config/containers/logging/splunk/
func initSplunkOpts() {
	pflag.String(splunk.TokenKey, "", "Splunk HTTP Event Collector token.")
	pflag.String(splunk.URLKey, "", "Path to your Splunk Enterprise, self-service Splunk Cloud instance, or Splunk Cloud managed cluster (including port and scheme used by HTTP Event Collector).")
	pflag.String(splunk.SourceKey, "", "Event source.")
	pflag.String(splunk.SourcetypeKey, "", "Event source type.")
	pflag.String(splunk.IndexKey, "", "Event index.")
	pflag.String(splunk.CapathKey, "", "Path to root certificate.")
	pflag.String(splunk.CanameKey, "", "Name to use for validating server certificate; by default the hostname of the splunk-url is used.")
	pflag.String(splunk.InsecureskipverifyKey, "", "Ignore server certificate validation.")
	pflag.String(splunk.FormatKey, "", "Message format. Can be inline, json or raw. Defaults to inline.")
	pflag.String(splunk.VerifyConnectionKey, "", "Verify on start, that docker can connect to Splunk server. Defaults to true.")
	pflag.String(splunk.GzipKey, "", "Enable/disable gzip compression to send events to Splunk Enterprise or Splunk Cloud instance. Defaults to false.")
	pflag.String(splunk.GzipLevelKey, "", "Set compression level for gzip. Valid values are -1 (default), 0 (no compression), 1 (best speed) ... 9 (best compression). Defaults to DefaultCompression.")
	pflag.String(splunk.SplunkTagKey, "", "Specify tag for message, which interpret some markup.")
	pflag.String(splunk.LabelsKey, "", "Comma-separated list of keys of labels, which should be included in message, if these labels are specified for container.")
	pflag.String(splunk.EnvKey, "", "Comma-separated list of keys of environment variables, which should be included in message, if these variables are specified for container.")
	pflag.String(splunk.EnvRegexKey, "", "Similar to and compatible with env. A regular expression to match logging-related environment variables. Used for advanced log tag options.")
}
