// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build !linux && !windows

package debug

import (
	"errors"
	"time"
)

// Not implemented.
func SendEventsToLog(_, _, _ string, _ time.Duration) {}

// Not implemented.
func StartStackTraceHandler() {}

// Not implemented.
func SetLogFilePath(_, _ string) error { return errors.New("not implemented") }

func FlushLog() {}
