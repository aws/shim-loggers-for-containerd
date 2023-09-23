// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build !windows
// +build !windows

// Package debug provides utilities for enhanced logging within the shim loggers for containerd, focusing on system
// journal integration and stack trace capture on signal reception.
package debug

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/journal"
)

var (
	// journalPriority is a map that maps strings to journal priority values.
	journalPriority = map[string]journal.Priority{
		ERROR: journal.PriErr,
		INFO:  journal.PriInfo,
		DEBUG: journal.PriDebug,
	}
)

// FlushLog only used for Windows file logging.
func FlushLog() {}

// SendEventsToLog dispatches log messages to the system journal based on the given priority.
func SendEventsToLog(syslogIdentifier string, msg string, msgType string, delay time.Duration) {
	journalType := journalPriority[msgType]
	sendEventsToJournal(syslogIdentifier, msg, journalType, delay)
}

// StartStackTraceHandler is used when the process catches signals, we will print the stack trace and write
// to system journal. This function is useful when debugging.
func StartStackTraceHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go func() {
		for range c {
			stackDump := make([]byte, 640*1024)
			n := runtime.Stack(stackDump, true)
			stackDump = stackDump[:n]
			SendEventsToLog(daemonName, fmt.Sprintf(
				"\n====== STACKTRACE ======\n%v\n%s\n====== /STACKTRACE ======\n",
				time.Now(),
				stackDump,
			), DEBUG, 2)
		}

		os.Exit(1)
	}()
}

// This is a temporary solution for logging the shim-loggers-for-containerd package itself. We directly
// send the events to system journal and they are identified by the package name. Since this process is
// started by containerd, we can check the logs using `journalctl -u containerd.service`.
func sendEventsToJournal(syslogIdentifier string, msg string, msgType journal.Priority, delay time.Duration) {
	vars := map[string]string{"SYSLOG_IDENTIFIER": syslogIdentifier}
	journal.Send(msg, msgType, vars) //nolint:errcheck,gosec // asynchronous process
	time.Sleep(delay * time.Second)
}

// SetLogFilePath only supported on Windows
// For non-Windows logs will be written to journald.
func SetLogFilePath(_, _ string) error {
	return errors.New("debugging to file not supported, debug logs will be written with journald")
}
