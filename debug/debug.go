// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package debug

const (
	daemonName = "shim-loggers-for-containerd"
	INFO       = "info"
	ERROR      = "err"
	DEBUG      = "debug"
)

var (
	// When this set to true, logger will print more events for debugging
	Verbose   = false
	LoggerErr error
)

func DeferFuncForRunLogDriver() {
	if LoggerErr != nil {
		SendEventsToLog(daemonName, LoggerErr.Error(), ERROR, 1)
	}
}
