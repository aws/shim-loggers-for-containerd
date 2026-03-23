// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build unit

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFetchFromEndpoint tests the fetchFromEndpoint helper with various HTTP
// response scenarios.
func TestFetchFromEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		expectErr    bool
		errContains  string
		expectedBody string
	}{
		{
			name:         "successful fetch returns body",
			statusCode:   http.StatusOK,
			responseBody: `{"token":"secret-value"}`,
			expectErr:    false,
			expectedBody: `{"token":"secret-value"}`,
		},
		{
			name:         "non-200 response returns error with status code and body",
			statusCode:   http.StatusNotFound,
			responseBody: `{"code":"NotFound","message":"not found"}`,
			expectErr:    true,
			errContains:  "HTTP 404",
		},
		{
			name:         "internal server error returns error with status code and body",
			statusCode:   http.StatusInternalServerError,
			responseBody: `{"code":"InternalServerError","message":"error"}`,
			expectErr:    true,
			errContains:  "HTTP 500",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = fmt.Fprint(w, tc.responseBody)
			}))
			defer server.Close()

			body, err := fetchFromEndpoint(server.URL)
			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
				assert.Contains(t, err.Error(), server.URL)
				// Verify the error includes the response body snippet for debugging.
				if tc.responseBody != "" {
					assert.Contains(t, err.Error(), tc.responseBody)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedBody, string(body))
			}
		})
	}
}

// TestFetchFromEndpointNetworkError tests that a network error returns a
// descriptive error including the URL.
func TestFetchFromEndpointNetworkError(t *testing.T) {
	t.Parallel()

	body, err := fetchFromEndpoint("http://127.0.0.1:0/nonexistent")
	require.Error(t, err)
	assert.Nil(t, body)
	assert.Contains(t, err.Error(), "failed to fetch from")
}
