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
		DatetimeFormat:      testDatetimeFormat,
		Endpoint:            testEndpoint,
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
		DatetimeFormatKey:      testDatetimeFormat,
		EndpointKey:            testEndpoint,
	}

	config, err := getAWSLogsConfig(args)
	require.NoError(t, err)
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
		EndpointKey:            testEndpoint,
	}

	config, err := getAWSLogsConfig(argsWithoutHeader)
	require.NoError(t, err)
	require.Equal(t, expectedConfig, config)
}

// TestGetAWSLogsConfigValidationError tests that getAWSLogsConfig returns an error when ValidateLogOpts fails.
func TestGetAWSLogsConfigValidationError(t *testing.T) {
	invalidArgs := &Args{
		Group:               testGroup,
		Region:              testRegion,
		Stream:              testStream,
		CredentialsEndpoint: testCredentialsEndpoint,
		CreateGroup:         testCreateGroup,
		// Specifying both 'awslogs-datetime-format' and 'awslogs-multiline-pattern' at the same time is an
		// invalid config.
		DatetimeFormat:   testDatetimeFormat,
		MultilinePattern: testMultilinePattern,
		Endpoint:         testEndpoint,
	}

	config, err := getAWSLogsConfig(invalidArgs)
	require.Error(t, err)
	require.Nil(t, config)
}

// TestLogFormatHeaderIsNotCompatibleWithDatetimeOrMultilineFormat tests if validateLogOptCompatability function
// correctly validates that LogFormatKey cannot be configured with DatetimeFormatKey or MultilineFormatKey.
func TestLogFormatHeaderIsNotCompatibleWithDatetimeOrMultilineFormat(t *testing.T) {
	_, err := getAWSLogsConfig(argsWithoutMultiline)
	require.EqualError(t, err,
		"you cannot configure log opt 'awslogs-datetime-format' or "+
			"'awslogs-multiline-pattern' when log opt 'awslogs-format' is set to 'json/emf'")
}

// TestValidateLogOptsDoesntErrorWithGoodConfig tests if validateLogOptCompatability function
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

	config, err := getAWSLogsConfig(argsValid)
	require.NoError(t, err)
	require.Equal(t, expectedConfig, config)
}

// TestGetAWSLogsConfigMinimal tests that getAWSLogsConfig correctly handles
// minimal configuration with only required fields.
func TestGetAWSLogsConfigMinimal(t *testing.T) {
	minimalArgs := &Args{
		Group:  testGroup,
		Region: testRegion,
		Stream: testStream,
	}

	expectedConfig := map[string]string{
		GroupKey:  testGroup,
		RegionKey: testRegion,
		StreamKey: testStream,
	}

	config, err := getAWSLogsConfig(minimalArgs)
	require.NoError(t, err)
	require.Equal(t, expectedConfig, config)
}
