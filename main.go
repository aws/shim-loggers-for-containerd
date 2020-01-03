package main

import (
	"os"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"
	"github.com/aws/shim-loggers-for-containerd/logger/awslogs"

	"github.com/containerd/containerd/runtime/v2/logging"
	"github.com/coreos/go-systemd/journal"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	awslogsDriverName = "awslogs"

	// mode and buffer size options
	modeKey          = "mode"
	maxBufferSizeKey = "max-buffer-size"

	// logDriver configuration
	containerIDKey   = "container-id"
	containerNameKey = "container-name"
	logDriverTypeKey = "log-driver"
)

func init() {
	initCommonLogOpts()
	initAWSLogsOpts()
}

func main() {
	pflag.Parse()
	if err := run(); err != nil {
		debug.SendEventsToJournal(logger.DaemonName, err.Error(), journal.PriErr)
		os.Exit(1)
	}
}

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

func run() error {
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return errors.Wrap(err, "unable to bind command line flags")
	}

	globalArgs, err := getGlobalArgs()
	if err != nil {
		return errors.Wrap(err, "unable to get global arguments")
	}

	logDriver := globalArgs.LogDriver
	debug.SendEventsToJournal(logger.DaemonName, "Driver: "+logDriver, journal.PriInfo)
	switch logDriver {
	case awslogsDriverName:
		awslogsArgs, err := getAWSLogsArgs()
		if err != nil {
			return errors.Wrap(err, "unable to get awslogs specified arguments")
		}
		loggerArgs := awslogs.InitLogger(globalArgs, awslogsArgs)
		logging.Run(loggerArgs.RunLogDriver)
	default:
		return errors.Errorf("unknown log driver: %s", logDriver)
	}

	return nil
}
