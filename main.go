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

package main

import (
	"os"
	"runtime"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"
	"github.com/aws/shim-loggers-for-containerd/logger/awslogs"
	"github.com/aws/shim-loggers-for-containerd/logger/fluentd"
	"github.com/aws/shim-loggers-for-containerd/logger/splunk"

	"github.com/containerd/containerd/runtime/v2/logging"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	initCommonLogOpts()
	initWindowsOpts()
	initDockerConfigOpts()
	initAWSLogsOpts()
	initFluentdOpts()
	initSplunkOpts()
}

func main() {
	pflag.Parse()
	if err := run(); err != nil {
		debug.SendEventsToLog(logger.DaemonName, err.Error(), debug.ERROR, 1)
		os.Exit(1)
	}
}

func run() error {
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return errors.Wrap(err, "unable to bind command line flags")
	}

	globalArgs, err := getGlobalArgs()
	if err != nil {
		return errors.Wrap(err, "unable to get global arguments")
	}

	// Read the Windows specific options and set the environment up accordingly
	if runtime.GOOS == "windows" {
		windowsArgs := getWindowsArgs()
		err = setWindowsEnv(windowsArgs.LogFileDir, globalArgs.ContainerName, windowsArgs.ProxyEnvVar)
		if err != nil {
			return errors.Wrap(err, "failed to set up Windows env with options")
		}
		defer cleanWindowsEnv(windowsArgs.ProxyEnvVar)
	}

	debug.Verbose = viper.GetBool(verboseKey)
	if debug.Verbose {
		debug.SendEventsToLog(logger.DaemonName, "Using verbose mode", debug.INFO, 0)
		// If in Verbose mode, start a goroutine to catch os signal and print stack trace
		debug.StartStackTraceHandler()
	}
	// Set UID and/or GID of main goroutine/shim logger process if specified.
	// If you are building with go version includes the following commit, you only need
	// to call this once in main goroutine. Otherwise you need call this function in all
	// goroutines to let this syscall work properly.
	// Commit: https://github.com/golang/go/commit/d1b1145cace8b968307f9311ff611e4bb810710c
	// TODO: remove the above comment once the changes are released: https://go-review.googlesource.com/c/go/+/210639
	if err = logger.SetUIDAndGID(globalArgs.UID, globalArgs.GID); err != nil {
		return err
	}

	logDriver := globalArgs.LogDriver
	debug.SendEventsToLog(logger.DaemonName, "Driver: "+logDriver, debug.INFO, 0)
	switch logDriver {
	case awslogsDriverName:
		if err := runAWSLogsDriver(globalArgs); err != nil {
			return errors.Wrap(err, "unable to run awslogs driver")
		}
	case fluentdDriverName:
		runFluentdDriver(globalArgs)
	case splunkDriverName:
		if err := runSplunkDriver(globalArgs); err != nil {
			return errors.Wrap(err, "unable to run splunk driver")
		}
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

func runSplunkDriver(globalArgs *logger.GlobalArgs) error {
	dockerConfigs, err := getDockerConfigs()
	if err != nil {
		return errors.Wrap(err, "unable to get docker config arguments")
	}

	args, err := getSplunkArgs()
	if err != nil {
		return errors.Wrap(err, "unable to get splunk specified arguments")
	}

	loggerArgs := splunk.InitLogger(globalArgs, dockerConfigs, args)
	logging.Run(loggerArgs.RunLogDriver)

	return nil
}

// setWindowsEnv reads the Windows options and sets them up
func setWindowsEnv(logDir, containerName, proxyEnvVar string) error {
	if logDir != "" {
		err := debug.SetLogFilePath(logDir, containerName)
		if err != nil {
			// Will include an error line if log-file-dir option is set for non-Windows in logs
			// Will ignore and continue to log with journald
			debug.SendEventsToLog(logger.DaemonName, err.Error(), debug.ERROR, 1)
			return err
		}
	}
	// proxyEnvVar will set the HTTP_PROXY and HTTPS_PROXY environment variables
	if proxyEnvVar != "" {
		err := os.Setenv("HTTP_PROXY", proxyEnvVar)
		if err != nil {
			return err
		}
		err = os.Setenv("HTTPS_PROXY", proxyEnvVar)
		if err != nil {
			return err
		}
	}
	return nil
}

// cleanWindowsEnv flushes the file logs for Windows and unsets the proxy env variables
func cleanWindowsEnv(proxyEnvVar string) {
	debug.FlushLog()
	if proxyEnvVar != "" {
		os.Unsetenv("HTTP_PROXY")
		os.Unsetenv("HTTPS_PROXY")
	}
}
