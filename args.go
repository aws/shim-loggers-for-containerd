// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"
	"github.com/aws/shim-loggers-for-containerd/logger/awslogs"
	"github.com/aws/shim-loggers-for-containerd/logger/fluentd"
	"github.com/aws/shim-loggers-for-containerd/logger/splunk"

	"github.com/docker/go-units"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultMaxBufferSize = "1m"
	defaultCleanupTime   = "5s"
	blockingMode         = "blocking"
	nonBlockingMode      = "non-blocking"
)

// getGlobalArgs get arguments that used for any log drivers
func getGlobalArgs() (*logger.GlobalArgs, error) {
	containerID, err := getRequiredValue(containerIDKey)
	if err != nil {
		return nil, err
	}
	containerName, err := getRequiredValue(containerNameKey)
	if err != nil {
		return nil, err
	}
	logDriver, err := getRequiredValue(logDriverTypeKey)
	if err != nil {
		return nil, err
	}
	mode, maxBufferSize, err := getModeAndMaxBufferSize()
	if err != nil {
		return nil, fmt.Errorf("unable to get value of flag %s and %s: %w", modeKey, maxBufferSizeKey, err)
	}
	cleanupTime, err := getCleanupTime()
	if err != nil {
		return nil, err
	}

	if debug.Verbose {
		debug.SendEventsToLog(logger.DaemonName,
			fmt.Sprintf("Container ID: %s, Container Name: %s, log driver: %s, mode: %s, max buffer size: %d",
				containerID, containerName, logDriver, mode, maxBufferSize),
			debug.DEBUG, 0)
	}

	args := &logger.GlobalArgs{
		ContainerID:   containerID,
		ContainerName: containerName,
		LogDriver:     logDriver,
		Mode:          mode,
		MaxBufferSize: maxBufferSize,
		UID:           viper.GetInt(uidKey),
		GID:           viper.GetInt(gidKey),
		CleanupTime:   cleanupTime,
	}

	return args, nil
}

// getWindowsArgs gets the optional Windows arguments
func getWindowsArgs() *logger.WindowsArgs {
	proxyVar := viper.GetString(ProxyEnvVarKey)
	logDir := viper.GetString(LogFileDirKey)

	return &logger.WindowsArgs{
		ProxyEnvVar: proxyVar,
		LogFileDir:  logDir,
	}
}

// getDockerConfigs gets the optional docker config variables
func getDockerConfigs() (*logger.DockerConfigs, error) {
	containerLabelsString := viper.GetString(ContainerLabelsKey)
	containerLabels := make(map[string]string)
	if containerLabelsString != "" {
		err := json.Unmarshal([]byte(containerLabelsString), &containerLabels)
		if err != nil {
			return nil, err
		}
	}

	containerEnvString := viper.GetString(ContainerEnvKey)
	containerEnv := make([]string, 0)
	if containerEnvString != "" {
		containerEnvMap := make(map[string]string)
		err := json.Unmarshal([]byte(containerEnvString), &containerEnvMap)
		if err != nil {
			return nil, err
		}
		// Docker logging use a slice to store the environment variables
		// Each item is in the format of "key=value"
		// ref: https://github.com/moby/moby/blob/c833222d54c00d64a0fc44c561a5973ecd414053/daemon/logger/loginfo.go#L60
		for envKey, envVal := range containerEnvMap {
			containerEnv = append(containerEnv, envKey+"="+envVal)
		}
	}

	// Docker config variables are optional
	// If no docker config variables are passed to shim logger, we will use empty values for them.
	args := &logger.DockerConfigs{
		ContainerImageID:   viper.GetString(ContainerImageIDKey),
		ContainerImageName: viper.GetString(ContainerImageNameKey),
		ContainerEnv:       containerEnv,
		ContainerLabels:    containerLabels,
	}

	return args, nil
}

// getAWSLogsArgs gets awslogs specified arguments for awslogs driver
func getAWSLogsArgs() (*awslogs.Args, error) {
	group, err := getRequiredValue(awslogs.GroupKey)
	if err != nil {
		return nil, err
	}
	region, err := getRequiredValue(awslogs.RegionKey)
	if err != nil {
		return nil, err
	}
	stream, err := getRequiredValue(awslogs.StreamKey)
	if err != nil {
		return nil, err
	}
	credentialsEndpoint, err := getRequiredValue(awslogs.CredentialsEndpointKey)
	if err != nil {
		return nil, err
	}

	return &awslogs.Args{
		Group:               group,
		Region:              region,
		Stream:              stream,
		CredentialsEndpoint: credentialsEndpoint,
		CreateGroup:         viper.GetString(awslogs.CreateGroupKey),
		CreateStream:        viper.GetString(awslogs.CreateStreamKey),
		MultilinePattern:    viper.GetString(awslogs.MultilinePatternKey),
		DatetimeFormat:      viper.GetString(awslogs.DatetimeFormatKey),
		Endpoint:            viper.GetString(awslogs.EndpointKey),
	}, nil
}

// getFluentdArgs gets fluentd specified arguments for fluentd log driver
func getFluentdArgs() *fluentd.Args {
	address := viper.GetString(fluentd.AddressKey)
	tag := viper.GetString(fluentd.FluentdTagKey)

	ac := viper.GetBool(fluentd.AsyncConnectKey)
	asyncConnect := strconv.FormatBool(ac)

	precision := viper.GetBool(fluentd.SubsecondPrecisionKey)
	subsecondPrecision := strconv.FormatBool(precision)

	buffer := viper.GetInt(fluentd.BufferLimitKey)
	bufferLimit := ""
	if buffer > 0 {
		bufferLimit = strconv.Itoa(buffer)
	}

	return &fluentd.Args{
		Address:            address,
		Tag:                tag,
		AsyncConnect:       asyncConnect,
		SubsecondPrecision: subsecondPrecision,
		BufferLimit:        bufferLimit,
	}
}

// getSplunkArgs gets Splunk specified arguments for Splunk log driver
func getSplunkArgs() (*splunk.Args, error) {
	token, err := getRequiredValue(splunk.TokenKey)
	if err != nil {
		return nil, err
	}
	url, err := getRequiredValue(splunk.URLKey)
	if err != nil {
		return nil, err
	}

	return &splunk.Args{
		Token:              token,
		URL:                url,
		Source:             viper.GetString(splunk.SourceKey),
		Sourcetype:         viper.GetString(splunk.SourcetypeKey),
		Index:              viper.GetString(splunk.IndexKey),
		Capath:             viper.GetString(splunk.CapathKey),
		Caname:             viper.GetString(splunk.CanameKey),
		Insecureskipverify: viper.GetString(splunk.InsecureskipverifyKey),
		Format:             viper.GetString(splunk.FormatKey),
		VerifyConnection:   viper.GetString(splunk.VerifyConnectionKey),
		Gzip:               viper.GetString(splunk.GzipKey),
		GzipLevel:          viper.GetString(splunk.GzipLevelKey),
		Tag:                viper.GetString(splunk.SplunkTagKey),
		TagSpecified:       isFlagPassed(splunk.SplunkTagKey),
		Labels:             viper.GetString(splunk.LabelsKey),
		Env:                viper.GetString(splunk.EnvKey),
		EnvRegex:           viper.GetString(splunk.EnvRegexKey),
	}, nil
}

// getRequiredValue parses required arguments or exits if any is missing
func getRequiredValue(flag string) (string, error) {
	isSet := viper.IsSet(flag)
	if !isSet {
		return "", fmt.Errorf("%s is required", flag)
	}
	val := viper.GetString(flag)

	return val, nil
}

// getModeAndMaxBufferSize gets mode option and max buffer size if in blocking mode
func getModeAndMaxBufferSize() (string, int, error) {
	var (
		mode       string
		maxBufSize int
		err        error
	)

	mode = viper.GetString(modeKey)
	switch mode {
	case "":
		mode = blockingMode
	case blockingMode:
	case nonBlockingMode:
		maxBufSize, err = getMaxBufferSize()
		if err != nil {
			return "", 0, fmt.Errorf("unable to get max buffer size: %w", err)
		}
	default:
		return "", 0, fmt.Errorf("unknown mode type: %s", mode)
	}

	return mode, maxBufSize, nil
}

// getMaxBufferSize gets either customer asked buffer size or default size 1m
func getMaxBufferSize() (int, error) {
	var (
		size int64
		err  error
	)
	maxBufferSize := viper.GetString(maxBufferSizeKey)
	if maxBufferSize == "" {
		size, err = units.RAMInBytes(defaultMaxBufferSize)
	} else {
		size, err = units.RAMInBytes(maxBufferSize)
	}

	if err != nil {
		return 0, fmt.Errorf("unable to parse buffer size to bytes: %w", err)
	}

	return int(size), nil
}

// getCleanupTime gets either customized cleanup time or default duration of 5s
func getCleanupTime() (*time.Duration, error) {
	cleanupTime := viper.GetString(cleanupTimeKey)
	if cleanupTime == "" {
		cleanupTime = defaultCleanupTime
	}
	duration, err := time.ParseDuration(cleanupTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse clean up time: %w", err)
	}
	if duration > time.Duration(12*time.Second) {
		return nil, fmt.Errorf("invalid time %s, maximum timeout is 12 seconds", duration.String())
	}

	return &duration, nil
}

// isFlagPassed determines whether a flag was passed by the client.
func isFlagPassed(name string) bool {
	passed := false
	pflag.Visit(func(f *pflag.Flag) {
		if f.Name == name {
			passed = true
		}
	})
	return passed
}
