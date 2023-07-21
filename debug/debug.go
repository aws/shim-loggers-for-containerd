// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package debug

const (
	daemonName = "shim-loggers-for-containerd"
	// INFO represents a log level for informational messages.
	INFO = "info"
	// ERROR represents a log level for error messages.
	ERROR = "err"
	// DEBUG represents a log level for debugging messages.
	DEBUG = "debug"
)

var (
	// Verbose indicates if additional debug events should be logged.
	Verbose = false
	// ErrLogger holds any errors related to the logger setup or execution.
	ErrLogger error
)

// DeferFuncForRunLogDriver checks and sends logger errors to the system log.
func DeferFuncForRunLogDriver() {
	if ErrLogger != nil {
		SendEventsToLog(daemonName, ErrLogger.Error(), ERROR, 1)
	}
}
