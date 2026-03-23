// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// SplunkTokenResponse is the JSON response body returned by the splunk token endpoint.
type SplunkTokenResponse struct {
	Token string `json:"token"`
}

// ContainerEnvResponse is the JSON response body returned by the container env endpoint.
type ContainerEnvResponse struct {
	Env map[string]string `json:"env"`
}

// httpClient is used by fetchFromEndpoint so that requests have a bounded deadline.
var httpClient = &http.Client{Timeout: 5 * time.Second}

// maxResponseBytes is the upper bound on how many bytes fetchFromEndpoint will
// read from any endpoint response. This prevents unbounded memory consumption
// if an endpoint misbehaves and returns a very large body.
const maxResponseBytes = 1 << 20 // 1 MiB

// maxErrBodyBytes is the upper bound on how many bytes of an error response
// body are included in the returned error message.
const maxErrBodyBytes = 512

// fetchFromEndpoint issues an HTTP GET to the given URL and returns the response
// body bytes. On non-200 status or network error, it returns a descriptive error
// including the HTTP status code, endpoint URL, and a truncated snippet of the
// response body.
func fetchFromEndpoint(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from %s: %w", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrBodyBytes))
		return nil, fmt.Errorf("failed to fetch from %s: HTTP %d: %s", url, resp.StatusCode, string(errBody))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %s: %w", url, err)
	}
	return body, nil
}
