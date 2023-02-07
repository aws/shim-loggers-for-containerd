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
	testCredentialsEndpoint = "test-credential-endpoints"
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
