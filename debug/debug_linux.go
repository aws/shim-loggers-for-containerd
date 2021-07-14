// +build !windows

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

package debug

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/journal"
	"github.com/pkg/errors"
)

var (
	// journalPriority is a map that maps strings to journal priority values
	journalPriority = map[string]journal.Priority{
		ERROR	  :	journal.PriErr,
		INFO	  :	journal.PriInfo,
		DEBUG	  :	journal.PriDebug,
	}
)

// FlushLog only used for Windows file logging
func FlushLog() {}

func SendEventsToLog(syslogIdentifier string, msg string, msgType string, delaySeconds time.Duration) {
	journalType := journalPriority[msgType]
	sendEventsToJournal(syslogIdentifier, msg, journalType, delaySeconds)
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
func sendEventsToJournal(syslogIdentifier string, msg string, msgType journal.Priority, delaySeconds time.Duration) {
	vars := map[string]string{"SYSLOG_IDENTIFIER": syslogIdentifier}
	journal.Send(msg, msgType, vars) //nolint:errcheck
	time.Sleep(delaySeconds * time.Second)
}

// SetLogFilePath only supported on Windows
// For non-Windows logs will be written to journald
func SetLogFilePath(logFlag, cId string) error{
	return errors.New("debugging to file not supported, debug logs will be written with journald")
}
