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
	}

	config := getAWSLogsConfig(args)
	require.Equal(t, expectedConfig, config)
}
