package main

import (
	"os"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"
	"github.com/aws/shim-loggers-for-containerd/logger/awslogs"
	"github.com/aws/shim-loggers-for-containerd/logger/fluentd"

	"github.com/containerd/containerd/runtime/v2/logging"
	"github.com/coreos/go-systemd/journal"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	initCommonLogOpts()
	initAWSLogsOpts()
	initFluentdOpts()
}

func main() {
	pflag.Parse()
	if err := run(); err != nil {
		debug.SendEventsToJournal(logger.DaemonName, err.Error(), journal.PriErr)
		os.Exit(1)
	}
}

func run() error {
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return errors.Wrap(err, "unable to bind command line flags")
	}

	debug.Verbose = viper.GetBool(verboseKey)
	if debug.Verbose {
		debug.StartStackTraceHandler()
	}

	globalArgs, err := getGlobalArgs()
	if err != nil {
		return errors.Wrap(err, "unable to get global arguments")
	}

	logDriver := globalArgs.LogDriver
	debug.SendEventsToJournal(logger.DaemonName, "Driver: "+logDriver, journal.PriInfo)
	switch logDriver {
	case awslogsDriverName:
		if err := runAWSLogsDriver(globalArgs); err != nil {
			return errors.Wrap(err, "unable to run awslogs driver")
		}
	case fluentdDriverName:
		runFluentdDriver(globalArgs)
	default:
		return errors.Errorf("unknown log driver: %s", logDriver)
	}

	return nil
}

func runAWSLogsDriver(globalArgs *logger.GlobalArgs) error {
	args, err := getAWSLogsArgs()
	if err != nil {
		return errors.Wrap(err, "unable to get awslogs specified arguments")
	}
	loggerArgs := awslogs.InitLogger(globalArgs, args)
	logging.Run(loggerArgs.RunLogDriver)

	return nil
}

func runFluentdDriver(globalArgs *logger.GlobalArgs) {
	args := getFluentdArgs()
	loggerArgs := fluentd.InitLogger(globalArgs, args)
	logging.Run(loggerArgs.RunLogDriver)
}
