// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build unit
// +build unit

package awslogs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testGroup               = "test-group"
	testRegion              = "test-region"
	testStream              = "test-stream"
	testCredentialsEndpoint = "test-credential-endpoints" //nolint:gosec // not credentials
	testCreateGroup         = "true"
	testCreateStream        = "true"
	testMultilinePattern    = "test-multiline-pattern"
	testDatetimeFormat      = "test-datetime-format"
	testEndpoint            = "test-endpoint"
	testJSONEmfLogformat    = "json/emf"
)

var (
	args = &Args{
		Group:               testGroup,
		Region:              testRegion,
		Stream:              testStream,
		CredentialsEndpoint: testCredentialsEndpoint,
		CreateGroup:         testCreateGroup,
		CreateStream:        testCreateStream,
		MultilinePattern:    testMultilinePattern,
		DatetimeFormat:      testDatetimeFormat,
		Endpoint:            testEndpoint,
		LogsFormatHeader:    testJSONEmfLogformat,
	}
)

var (
	argsWithoutHeader = &Args{
		Group:               testGroup,
		Region:              testRegion,
		Stream:              testStream,
		CredentialsEndpoint: testCredentialsEndpoint,
		CreateGroup:         testCreateGroup,
		CreateStream:        testCreateStream,
		MultilinePattern:    testMultilinePattern,
		DatetimeFormat:      testDatetimeFormat,
		Endpoint:            testEndpoint,
	}
)

var (
	argsWithoutMultiline = &Args{
		Group:               testGroup,
		Region:              testRegion,
		Stream:              testStream,
		CredentialsEndpoint: testCredentialsEndpoint,
		CreateGroup:         testCreateGroup,
		CreateStream:        testCreateStream,
		DatetimeFormat:      testDatetimeFormat,
		Endpoint:            testEndpoint,
		LogsFormatHeader:    testJSONEmfLogformat,
	}
)

var (
	argsValid = &Args{
		Group:               testGroup,
		Region:              testRegion,
		Stream:              testStream,
		CredentialsEndpoint: testCredentialsEndpoint,
		CreateGroup:         testCreateGroup,
		CreateStream:        testCreateStream,
		MultilinePattern:    testMultilinePattern,
		Endpoint:            testEndpoint,
	}
)

// TestGetAWSLogsConfig tests if getAWSLogsConfig function converts log config
// maps correctly.
func TestGetAWSLogsConfig(t *testing.T) {
	expectedConfig := map[string]string{
		GroupKey:               testGroup,
		RegionKey:              testRegion,
		StreamKey:              testStream,
		CredentialsEndpointKey: testCredentialsEndpoint,
		CreateGroupKey:         testCreateGroup,
		CreateStreamKey:        testCreateStream,
		MultilinePatternKey:    testMultilinePattern,
		DatetimeFormatKey:      testDatetimeFormat,
		EndpointKey:            testEndpoint,
		LogFormatKey:           testJSONEmfLogformat,
	}

	config := getAWSLogsConfig(args)
	require.Equal(t, expectedConfig, config)
}

// TestGetAWSLogsConfigWithoutFormatHeader tests if getAWSLogsConfig function converts log config
// maps correctly if LogsFormatHeader is omitted.
func TestGetAWSLogsConfigWithoutFormatHeader(t *testing.T) {
	expectedConfig := map[string]string{
		GroupKey:               testGroup,
		RegionKey:              testRegion,
		StreamKey:              testStream,
		CredentialsEndpointKey: testCredentialsEndpoint,
		CreateGroupKey:         testCreateGroup,
		CreateStreamKey:        testCreateStream,
		MultilinePatternKey:    testMultilinePattern,
		DatetimeFormatKey:      testDatetimeFormat,
		EndpointKey:            testEndpoint,
	}

	config := getAWSLogsConfig(argsWithoutHeader)
	require.Equal(t, expectedConfig, config)
}

// TestLogFormatHeaderIsNotCompatibleWithDatetimeOrMultilineFormat tests if validateLogOptCompatability function
// correctly validates that LogFormatKey cannot be configured with DatetimeFormatKey or MultilineFormatKey.
func TestLogFormatHeaderIsNotCompatibleWithDatetimeOrMultilineFormat(t *testing.T) {
	expectedConfig := map[string]string{
		GroupKey:               testGroup,
		RegionKey:              testRegion,
		StreamKey:              testStream,
		CredentialsEndpointKey: testCredentialsEndpoint,
		CreateGroupKey:         testCreateGroup,
		CreateStreamKey:        testCreateStream,
		DatetimeFormatKey:      testDatetimeFormat,
		EndpointKey:            testEndpoint,
		LogFormatKey:           testJSONEmfLogformat,
	}

	config := getAWSLogsConfig(argsWithoutMultiline)
	require.Equal(t, expectedConfig, config)
	err := validateLogOptCompatability(config)
	require.EqualError(t, err,
		"you cannot configure log opt 'awslogs-datetime-format' or "+
			"'awslogs-multiline-pattern' when log opt 'awslogs-format' is set to 'json/emf'")
}

// TestLogFormatHeaderIsNotCompatibleWithDatetimeOrMultilineFormat tests if validateLogOptCompatability function
// correctly validates that LogFormatKey cannot be configured with DatetimeFormatKey or MultilineFormatKey.
func TestValidateLogOptsDoesntErrorWithGoodConfig(t *testing.T) {
	expectedConfig := map[string]string{
		GroupKey:               testGroup,
		RegionKey:              testRegion,
		StreamKey:              testStream,
		CredentialsEndpointKey: testCredentialsEndpoint,
		CreateGroupKey:         testCreateGroup,
		CreateStreamKey:        testCreateStream,
		MultilinePatternKey:    testMultilinePattern,
		EndpointKey:            testEndpoint,
	}

	config := getAWSLogsConfig(argsValid)
	require.Equal(t, expectedConfig, config)
	err := validateLogOptCompatability(config)
	require.NoError(t, err)
}
