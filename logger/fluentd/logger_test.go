// +build unit

package fluentd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testAddress      = "testAddress"
	testAsyncConnect = "false"
	testTag          = "testTag"
)

var (
	args = &Args{
		Address:      testAddress,
		AsyncConnect: testAsyncConnect,
		Tag:          testTag,
	}
)

func TestGetFluentdConfig(t *testing.T) {
	expectedConfig := map[string]string{
		AddressKey:      testAddress,
		AsyncConnectKey: testAsyncConnect,
		tagKey:          testTag,
	}

	config := getFluentdConfig(args)
	require.Equal(t, expectedConfig, config)
}
