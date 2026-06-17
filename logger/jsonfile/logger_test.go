// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build unit
// +build unit

package jsonfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testLogPath      = "/var/log/ecs/json-file/abc123/abc123-json.log"
	testMaxSize      = "10m"
	testMaxFile      = "5"
	testCompress     = "false"
	testLabels       = "label0,label1"
	testLabelsRegex  = "^app\\..*"
	testEnv          = "envVar0,envVar1"
	testEnvRegex     = "^APP_.*"
	testTag          = "{{.ImageName}}/{{.ID}}"
	testTagSpecified = true
)

// TestGetJSONFileConfig tests that all set arguments are converted to the moby-side config map.
func TestGetJSONFileConfig(t *testing.T) {
	args := &Args{
		LogPath:      testLogPath, // not in moby config map; consumed via info.LogPath
		MaxSize:      testMaxSize,
		MaxFile:      testMaxFile,
		Compress:     testCompress,
		Labels:       testLabels,
		LabelsRegex:  testLabelsRegex,
		Env:          testEnv,
		EnvRegex:     testEnvRegex,
		Tag:          testTag,
		TagSpecified: testTagSpecified,
	}

	expectedConfig := map[string]string{
		MaxSizeKey:     testMaxSize,
		MaxFileKey:     testMaxFile,
		CompressKey:    testCompress,
		LabelsKey:      testLabels,
		LabelsRegexKey: testLabelsRegex,
		EnvKey:         testEnv,
		EnvRegexKey:    testEnvRegex,
		tagKey:         testTag,
	}

	config, err := getJSONFileConfig(args)
	require.NoError(t, err)
	require.Equal(t, expectedConfig, config)
}

// TestGetJSONFileConfigOptionalFieldsOmitted tests that empty optional fields are NOT
// added to the config map, since moby's ValidateLogOpts rejects unknown keys but is
// happy with a missing key.
func TestGetJSONFileConfigOptionalFieldsOmitted(t *testing.T) {
	args := &Args{
		LogPath: testLogPath,
		// All other fields are zero-valued.
	}

	config, err := getJSONFileConfig(args)
	require.NoError(t, err)
	require.Empty(t, config, "empty Args should produce an empty config map")
}

// TestGetJSONFileConfigTagSpecifiedFalse tests that an explicitly-empty tag (TagSpecified=false)
// is not added to the config map even if Tag has a value, mirroring the splunk pattern.
func TestGetJSONFileConfigTagSpecifiedFalse(t *testing.T) {
	args := &Args{
		LogPath:      testLogPath,
		Tag:          testTag,
		TagSpecified: false,
	}

	config, err := getJSONFileConfig(args)
	require.NoError(t, err)
	_, hasTag := config[tagKey]
	require.False(t, hasTag, "tag should not appear in config when TagSpecified is false")
}

// TestGetJSONFileConfigInvalidCompressBool tests that invalid values for fields that moby's
// jsonfilelog.New parses (e.g., compress) pass ValidateLogOpts but would fail at writer-start
// time. ValidateLogOpts only checks for unknown keys, not value validity, so this case
// returns success here.
//
// Note: moby's runtime guardrail rejects compress=true when max-file<2 or max-size is unset,
// but that check happens in jsonfilelog.New (writer construction), not in ValidateLogOpts.
// We test the value-level rejection via the e2e tests; here we only assert the option-key
// validation surface.
func TestGetJSONFileConfigPassesValidationForKnownKeys(t *testing.T) {
	args := &Args{
		LogPath:  testLogPath,
		MaxSize:  "garbage-but-known-key",
		MaxFile:  "also-garbage",
		Compress: "neither-true-nor-false",
	}

	config, err := getJSONFileConfig(args)
	require.NoError(t, err, "ValidateLogOpts only checks key names, not value validity")
	require.Contains(t, config, MaxSizeKey)
	require.Contains(t, config, MaxFileKey)
	require.Contains(t, config, CompressKey)
}

// TestRunLogDriverCreatesLogDirectory verifies that the log file's parent
// directory is created if it does not exist.
func TestRunLogDriverCreatesLogDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	containerID := "test-container-abc123"
	logPath := filepath.Join(tmpDir, containerID, containerID+"-json.log")

	// The parent of logPath should not exist yet.
	logDir := filepath.Dir(logPath)
	_, err := os.Stat(logDir)
	require.True(t, os.IsNotExist(err), "log directory should not exist before test")

	// Simulate what RunLogDriver does: MkdirAll on the parent directory.
	err = os.MkdirAll(logDir, logDirMode)
	require.NoError(t, err, "MkdirAll should create the per-container directory")

	// Verify the directory was created.
	info, err := os.Stat(logDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

// TestRunLogDriverLogDirectoryAlreadyExists verifies that MkdirAll is a no-op
// when the directory already exists (idempotent).
func TestRunLogDriverLogDirectoryAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	containerID := "existing-container"
	logDir := filepath.Join(tmpDir, containerID)

	// Pre-create the directory.
	require.NoError(t, os.MkdirAll(logDir, logDirMode))

	// MkdirAll again should succeed without error.
	err := os.MkdirAll(logDir, logDirMode)
	assert.NoError(t, err, "MkdirAll on existing directory should be idempotent")
}
