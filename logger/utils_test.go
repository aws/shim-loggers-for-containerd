// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build unit
// +build unit

package logger

import (
	"errors"
	"testing"
	"time"
)

func TestRetryWithExponentialBackoff(t *testing.T) {
	t.Run("succeeds first try", func(t *testing.T) {
		attempts := 0
		operation := func() error {
			attempts++
			return nil
		}

		err := retryWithExponentialBackoff(
			operation,
			100*time.Millisecond,
			1*time.Second,
		)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if attempts != 1 {
			t.Errorf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("succeeds after retries", func(t *testing.T) {
		attempts := 0
		operation := func() error {
			attempts++
			if attempts < 3 {
				return errors.New("temporary error")
			}
			return nil
		}

		err := retryWithExponentialBackoff(
			operation,
			100*time.Millisecond,
			1*time.Second,
		)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("respects max backoff", func(t *testing.T) {
		attempts := 0
		startTime := time.Now()
		operation := func() error {
			attempts++
			if attempts < 4 {
				return errors.New("temporary error")
			}
			return nil
		}

		err := retryWithExponentialBackoff(
			operation,
			100*time.Millisecond,
			200*time.Millisecond,
		)

		duration := time.Since(startTime)
		// Calculate theoretical max duration:
		// 1st try: immediate
		// 1st retry: 100ms + 25% jitter -> 125ms
		// 2nd retry: 200ms + 25% jitter -> 250ms
		// 3rd retry: 200ms + 25% jitter -> 250ms
		// Total maximum expected: ~625ms. Add 75ms delta to the expectation.
		maxExpectedDuration := 700 * time.Millisecond

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if attempts != 4 {
			t.Errorf("expected 4 attempts, got %d", attempts)
		}
		if duration > maxExpectedDuration {
			t.Errorf("backoff took too long. Expected less than %v, got %v", maxExpectedDuration, duration)
		}
	})

	t.Run("handles zero initial backoff", func(t *testing.T) {
		attempts := 0
		startTime := time.Now()
		operation := func() error {
			attempts++
			if attempts < 3 {
				return errors.New("temporary error")
			}
			return nil
		}

		err := retryWithExponentialBackoff(
			operation,
			0,
			1*time.Second,
		)

		duration := time.Since(startTime)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
		if duration > 10*time.Millisecond { // Should be very quick with zero backoff
			t.Errorf("expected quick execution with zero backoff, took %v", duration)
		}
	})

	t.Run("succeeds after 100 tries", func(t *testing.T) {
		attempts := 0
		operation := func() error {
			attempts++
			if attempts >= 100 {
				return nil
			}
			return errors.New("temporary error")
		}

		err := retryWithExponentialBackoff(
			operation,
			10*time.Millisecond,
			30*time.Millisecond,
		)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if attempts != 100 {
			t.Errorf("expected 100 attempt, got %d", attempts)
		}
	})
}
