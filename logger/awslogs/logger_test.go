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
	testMultilinePattern    = "test-multiline-pattern"
	testDatetimeFormat      = "test-datetime-format"
)

var (
	args = &Args{
		Group:               testGroup,
		Region:              testRegion,
		Stream:              testStream,
		CredentialsEndpoint: testCredentialsEndpoint,
		CreateGroup:         testCreateGroup,
		MultilinePattern:    testMultilinePattern,
		DatetimeFormat:      testDatetimeFormat,
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
		MultilinePatternKey:    testMultilinePattern,
		DatetimeFormatKey:      testDatetimeFormat,
	}

	config := getAWSLogsConfig(args)
	require.Equal(t, expectedConfig, config)
}
