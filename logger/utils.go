// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"fmt"
	"math"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"
)

func retryWithExponentialBackoff(operation func() error, initialBackoff time.Duration, maxBackoff time.Duration) error {
	attempt := 0
	reachMaxBackoff := false
	for {
		err := operation()
		if err == nil {
			return nil
		}

		// Calculate backoff duration with exponential increase and jitter.
		backoff := float64(initialBackoff) * math.Pow(2, float64(attempt))
		if backoff > float64(maxBackoff) {
			backoff = float64(maxBackoff)
			reachMaxBackoff = true
		}

		jitter := 0.25 * backoff // Add 25% jitter.
		sleepDuration := time.Duration(backoff + jitter)

		debug.SendEventsToLog(DaemonName, fmt.Sprintf("Retrying the operation in %v", sleepDuration), debug.ERROR, 0)
		time.Sleep(sleepDuration)
		if !reachMaxBackoff {
			attempt++
		}
	}
}
