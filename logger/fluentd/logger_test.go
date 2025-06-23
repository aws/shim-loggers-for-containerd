// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build unit
// +build unit

package fluentd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testAddress            = "testAddress"
	testAsyncConnect       = "false"
	testTag                = "testTag"
	testSubsecondPrecision = "true"
	testBufferLimit        = "12345"
	testWriteTimeout       = "1s"
)

var (
	args = &Args{
		Address:            testAddress,
		AsyncConnect:       testAsyncConnect,
		SubsecondPrecision: testSubsecondPrecision,
		Tag:                testTag,
		BufferLimit:        testBufferLimit,
		WriteTimeout:       testWriteTimeout,
	}
)

func TestGetFluentdConfig(t *testing.T) {
	expectedConfig := map[string]string{
		AddressKey:            testAddress,
		AsyncConnectKey:       testAsyncConnect,
		SubsecondPrecisionKey: testSubsecondPrecision,
		tagKey:                testTag,
		BufferLimitKey:        testBufferLimit,
		WriteTimeoutKey:       testWriteTimeout,
	}

	config := getFluentdConfig(args)
	require.Equal(t, expectedConfig, config)
}
