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
	"fmt"
	"strconv"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"
	"github.com/aws/shim-loggers-for-containerd/logger/awslogs"
	"github.com/aws/shim-loggers-for-containerd/logger/fluentd"
	"github.com/aws/shim-loggers-for-containerd/logger/splunk"

	"github.com/coreos/go-systemd/journal"
	"github.com/docker/go-units"
	"github.com/pkg/errors"
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
		return nil, errors.Wrapf(err, "unable to get value of flag %s and %s", modeKey, maxBufferSizeKey)
	}
	cleanupTime, err := getCleanupTime()
	if err != nil {
		return nil, err
	}

	if debug.Verbose {
		debug.SendEventsToJournal(logger.DaemonName,
			fmt.Sprintf("Container ID: %s, Container Name: %s, log driver: %s, mode: %s, max buffer size: %d",
				containerID, containerName, logDriver, mode, maxBufferSize),
			journal.PriDebug, 0)
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
		MultilinePattern:    viper.GetString(awslogs.MultilinePatternKey),
		DatetimeFormat:      viper.GetString(awslogs.DatetimeFormatKey),
	}, nil
}

// getFluentdArgs gets fluentd specified arguments for fluentd log driver
func getFluentdArgs() *fluentd.Args {
	address := viper.GetString(fluentd.AddressKey)
	tag := viper.GetString(fluentd.FluentdTagKey)

	ac := viper.GetBool(fluentd.AsyncConnectKey)
	asyncConnect := strconv.FormatBool(ac)

	return &fluentd.Args{
		Address:      address,
		Tag:          tag,
		AsyncConnect: asyncConnect,
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
		Labels:             viper.GetString(splunk.LabelsKey),
		Env:                viper.GetString(splunk.EnvKey),
		EnvRegex:           viper.GetString(splunk.EnvRegexKey),
	}, nil
}

// getRequiredValue parses required arguments or exits if any is missing
func getRequiredValue(flag string) (string, error) {
	isSet := viper.IsSet(flag)
	if !isSet {
		err := errors.Errorf("%s is required", flag)
		return "", err
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
			return "", 0, errors.Wrap(err, "unable to get max buffer size")
		}
	default:
		return "", 0, errors.Errorf("unknown mode type: %s", mode)
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
		return 0, errors.Wrap(err, "unable to parse buffer size to bytes")
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
		return nil, errors.Wrap(err, "failed to parse clean up time")
	}
	if duration > time.Duration(12*time.Second) {
		return nil, errors.Errorf("invalid time %s, maximum timeout is 12 seconds.", duration.String())
	}

	return &duration, nil
}
