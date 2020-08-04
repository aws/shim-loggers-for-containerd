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
)

var (
	args = &Args{
		Address:            testAddress,
		AsyncConnect:       testAsyncConnect,
		SubsecondPrecision: testSubsecondPrecision,
		Tag:                testTag,
	}
)

func TestGetFluentdConfig(t *testing.T) {
	expectedConfig := map[string]string{
		AddressKey:            testAddress,
		AsyncConnectKey:       testAsyncConnect,
		SubsecondPrecisionKey: testSubsecondPrecision,
		tagKey:                testTag,
	}

	config := getFluentdConfig(args)
	require.Equal(t, expectedConfig, config)
}
