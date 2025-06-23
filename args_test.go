// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build unit
// +build unit

package main

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aws/shim-loggers-for-containerd/logger/fluentd"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

const (
	testContainerID   = "test-container-id"
	testContainerName = "test-container-name"
	testLogDriver     = "test-log-driver"

	testContainerImageName = "test-container-image-name"
	testContainerImageID   = "test-container-image-id"
	testContainerLabels    = "{\"label0\":\"labelValue0\",\"label1\":\"labelValue1\"}"
	testContainerEnv       = "{\"env0\":\"envValue0\",\"env1\":\"envValue1\"}"
)

var (
	testContainerLabelsMap = map[string]string{
		"label0": "labelValue0",
		"label1": "labelValue1",
	}

	testContainerEnvSlice = [2]string{
		"env0=envValue0",
		"env1=envValue1",
	}
)

// TestGetGlobalArgs tests getGlobalArgs with/without correct settings of
// flag values.
func TestGetGlobalArgs(t *testing.T) {
	t.Run("NoError", testGetGlobalArgsNoError)
	t.Run("WithError", testGetGlobalArgsWithError)
}

// testGetGlobalArgsNoError is a sub-test of TestGetGlobalArgs. It tests
// getGlobalArgs function with no returned errors.
func testGetGlobalArgsNoError(t *testing.T) {
	// Unset all keys used for this test
	defer viper.Reset()

	viper.Set(containerIDKey, testContainerID)
	viper.Set(containerNameKey, testContainerName)
	viper.Set(logDriverTypeKey, testLogDriver)

	args, err := getGlobalArgs()
	require.NoError(t, err)
	assert.Equal(t, args.ContainerID, testContainerID)
	assert.Equal(t, args.ContainerName, testContainerName)
	assert.Equal(t, args.LogDriver, testLogDriver)
	assert.Equal(t, args.Mode, blockingMode)
	assert.Equal(t, args.MaxBufferSize, 0)
	assert.Equal(t, *args.CleanupTime, 5*time.Second)
}

// testGetGlobalArgsWithError is a sub-test of TestGetGlobalArgs. It tests
// getGlobalArgs function with returned errors when required flag values are
// missed. Since the logic in getGlobalArgs is examining required flag values
// one by one in this order:
//
//	container-id, container-name, log-driver,
//
// in this test, we will set flag one by one in the previous order to examine if
// each value check worked as expected.
func testGetGlobalArgsWithError(t *testing.T) {
	// Unset all keys used for this test
	defer viper.Reset()

	testCasesWithError := []struct {
		keyToSet string
		valToSet string
	}{
		{containerIDKey, testContainerID},
		{containerNameKey, testContainerName},
		{logDriverTypeKey, testLogDriver},
	}

	for _, tc := range testCasesWithError {
		_, err := getGlobalArgs()
		require.Error(t, err)
		require.Contains(t,
			err.Error(),
			fmt.Sprintf("%s is required", tc.keyToSet),
		)
		viper.Set(tc.keyToSet, tc.valToSet)
	}
}

// TestGetModeAndMaxBufferSize tests getModeAndMaxBufferSize with/without correct
// settings of mode.
func TestGetModeAndMaxBufferSize(t *testing.T) {
	t.Run("NoError", testGetModeAndMaxBufferSizeNoError)
	t.Run("WithError", testGetModeAndMaxBufferSizeWithError)
}

// testGetModeAndMaxBufferSizeNoError is a sub-test of TestGetModeAndMaxBufferSize.
// It tests different valid values set for mode option.
func testGetModeAndMaxBufferSizeNoError(t *testing.T) {
	// Unset all keys used for this test
	defer viper.Reset()

	testCasesNoError := []struct {
		mode               string
		expectedMode       string
		expectedBufferSize int
	}{
		{"", blockingMode, 0},
		{blockingMode, blockingMode, 0},
		{nonBlockingMode, nonBlockingMode, int(math.Pow(2, 20))},
	}

	for _, tc := range testCasesNoError {
		viper.Set(modeKey, tc.mode)
		mode, maxBufferSize, err := getModeAndMaxBufferSize()
		require.NoError(t, err)
		require.Equal(t, tc.expectedMode, mode)
		require.Equal(t, tc.expectedBufferSize, maxBufferSize)
	}
}

// testGetModeAndMaxBufferSizeNoError is a sub-test of TestGetModeAndMaxBufferSize. It
// tests invalid values set for mode option.
func testGetModeAndMaxBufferSizeWithError(t *testing.T) {
	// Unset all keys used for this test
	defer viper.Reset()

	viper.Set(modeKey, "test-mode")
	_, _, err := getModeAndMaxBufferSize()
	require.Error(t, err)
}

// TestGetMaxBufferSize tests getMaxBufferSize with/without valid setting max buffer
// size options.
func TestGetMaxBufferSize(t *testing.T) {
	t.Run("NoError", testGetMaxBufferSizeNoError)
	t.Run("WithError", testGetMaxBufferSizeWithError)
}

// testGetMaxBufferSizeNoError is a sub-test of TestGetMaxBufferSize. It tests
// getMaxBufferSize with multiple valid user-set values.
func testGetMaxBufferSizeNoError(t *testing.T) {
	// Unset all keys used for this test
	defer viper.Reset()

	testCasesNoError := []struct {
		bufferSize         string
		expectedBufferSize int
	}{
		{"", int(math.Pow(2, 20))},
		{"2m", int(math.Pow(2, 21))},
		{"4k", int(math.Pow(2, 12))},
		{"1234", 1234},
	}

	for _, tc := range testCasesNoError {
		viper.Set(maxBufferSizeKey, tc.bufferSize)
		size, err := getMaxBufferSize()
		require.NoError(t, err)
		require.Equal(t, tc.expectedBufferSize, size)
	}
}

// testGetMaxBufferSizeWithError is a sub-test of TestGetMaxBufferSize. It tests
// getMaxBufferSize with multiple invalid user-set values.
func testGetMaxBufferSizeWithError(t *testing.T) {
	// Unset all keys used for this test
	defer viper.Reset()

	testCasesWithError := []struct {
		bufferSize string
	}{
		{"3q"},
		{"-1"},
	}

	for _, tc := range testCasesWithError {
		viper.Set(maxBufferSizeKey, tc.bufferSize)
		_, err := getMaxBufferSize()
		require.Error(t, err)
	}
}

// TestGetCleanupTime tests getCleanupTime with/without valid setting cleanup time options.
func TestGetCleanupTime(t *testing.T) {
	t.Run("NoError", testGetCleanupTimeNoError)
	t.Run("WithError", testGetCleanupTimeWithError)
}

// testGetCleanupTimeNoError is a sub-test of TestGetCleanupTime. It tests getCleanupTime
// with multiple valid user-set values.
func testGetCleanupTimeNoError(t *testing.T) {
	// Unset all keys used for this test
	defer viper.Reset()

	testCasesNoError := []struct {
		cleanupTime         string
		expectedCleanupTime time.Duration
	}{
		{"3s", 3 * time.Second},
		{"10s", 10 * time.Second},
		{"12s", 12 * time.Second},
	}

	for _, tc := range testCasesNoError {
		viper.Set(cleanupTimeKey, tc.cleanupTime)
		cleanupTime, err := getCleanupTime()
		require.NoError(t, err)
		require.Equal(t, tc.expectedCleanupTime, *cleanupTime)
	}
}

// testGetCleanupTimeWithError is a sub-test of TestGetCleanupTime. It tests getCleanupTime
// with multiple invalid user-set values.
func testGetCleanupTimeWithError(t *testing.T) {
	// Unset all keys used for this test
	defer viper.Reset()

	testCasesWithError := []struct {
		cleanupTime string
	}{
		{"3"},
		{"15s"},
	}

	for _, tc := range testCasesWithError {
		viper.Set(cleanupTimeKey, tc.cleanupTime)
		_, err := getCleanupTime()
		require.Error(t, err)
	}
}

// TestGetDockerConfigs tests that we can correctly get the docker config input parameters.
func TestGetDockerConfigs(t *testing.T) {
	t.Run("NoError", testGetDockerConfigsNoError)
	t.Run("Empty", testGetDockerConfigsEmpty)
	t.Run("WithError", testGetDockerConfigsWithError)
}

// testGetDockerConfigsNoError tests that the correctly formatted input can be parsed without error.
func testGetDockerConfigsNoError(t *testing.T) {
	defer viper.Reset()

	viper.Set(ContainerImageNameKey, testContainerImageName)
	viper.Set(ContainerImageIDKey, testContainerImageID)
	viper.Set(ContainerLabelsKey, testContainerLabels)
	viper.Set(ContainerEnvKey, testContainerEnv)

	args, err := getDockerConfigs()
	require.NoError(t, err)
	assert.Equal(t, testContainerImageName, args.ContainerImageName)
	assert.Equal(t, testContainerImageID, args.ContainerImageID)
	assert.Equal(t, true, reflect.DeepEqual(args.ContainerLabels, testContainerLabelsMap))
	// Docker logging use a slice to store the environment variables
	// Each item is in the format of "key=value" the unit test should
	// convert it back to a map to guarantee content exists but disregard
	// the order.
	// ref: https://github.com/moby/moby/blob/c833222d54c00d64a0fc44c561a5973ecd414053/daemon/logger/loginfo.go#L60
	testContainerEnvMap := make(map[string]struct{})
	for _, v := range testContainerEnvSlice {
		testContainerEnvMap[v] = struct{}{}
	}
	argsContainerEnvMap := make(map[string]struct{})
	for _, v := range args.ContainerEnv {
		argsContainerEnvMap[v] = struct{}{}
	}
	assert.Equal(t, true, reflect.DeepEqual(testContainerEnvMap, argsContainerEnvMap))
}

// testGetDockerConfigsEmpty tests that empty docker config input parameter generates no error.
func testGetDockerConfigsEmpty(t *testing.T) {
	defer viper.Reset()

	_, err := getDockerConfigs()
	require.NoError(t, err)
}

// testGetDockerConfigsWithError tests that malformat docker config results in an error.
func testGetDockerConfigsWithError(t *testing.T) {
	defer viper.Reset()
	testCaseWithError := "{invalidJsonMap"

	viper.Set(ContainerLabelsKey, testCaseWithError)
	viper.Set(ContainerEnvKey, testCaseWithError)
	_, err := getDockerConfigs()
	require.Error(t, err)
}

// TestIsFlagPassed tests that we are correctly determining whether a flag is passed or not.
// Do not parallelize this test (or any test which interacts with flag parsing and OS args).
func TestIsFlagPassed(t *testing.T) {
	t.Run("No", testIsFlagPassedNo)
	t.Run("Yes", testIsFlagPassedYes)
}

// testIsFlagPassedNo creates a flag but doesn't pass it, and confirms that we don't interpret it to be set.
func testIsFlagPassedNo(t *testing.T) {
	defer func() {
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	}()

	pflag.String("test-flag", "", "test flag description")
	pflag.Parse()
	require.False(t, isFlagPassed("test-flag"))
}

// testIsFlagPassedYes creates a flag and mimics it being passed. Then, it confirms the flag was passed.
func testIsFlagPassedYes(t *testing.T) {
	defer func() {
		os.Args = os.Args[:len(os.Args)-1]
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	}()

	testFlag := "test-flag"
	pflag.String(testFlag, "", "test flag description")
	// pass the flag via os.Args
	os.Args = append(os.Args, "--"+testFlag+"=")
	pflag.Parse()
	require.True(t, isFlagPassed("test-flag"))
}

func TestGetFluentdArgs(t *testing.T) {
	for _, tc := range []struct {
		writeTimeoutValue string
		name              string
		expectedArgs      *fluentd.Args
	}{
		{
			writeTimeoutValue: "1s",
			name:              "overwrite works",
			expectedArgs: &fluentd.Args{
				AsyncConnect:       "false",
				SubsecondPrecision: "false",
				WriteTimeout:       "1s"},
		},
		{
			writeTimeoutValue: "",
			name:              "default works",
			expectedArgs: &fluentd.Args{
				AsyncConnect:       "false",
				SubsecondPrecision: "false",
				WriteTimeout:       "5s"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer viper.Reset()
			viper.Set(fluentd.WriteTimeoutKey, tc.writeTimeoutValue)
			args := getFluentdArgs()
			assert.DeepEqual(t, args, tc.expectedArgs)
		})
	}
}
