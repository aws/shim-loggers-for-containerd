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
)

const (
	daemonName = "shim-loggers-for-containerd"
)

var (
	// When this set to true, logger will print more events for debugging
	Verbose   = false
	LoggerErr error
)

// This is a temporary solution for logging the shim-loggers-for-containerd package itself. We directly
// send the events to system journal and they are identified by the package name. Since this process is
// started by containerd, we can check the logs using `journalctl -u containerd.service`.
func SendEventsToJournal(syslogIdentifier string, msg string, msgType journal.Priority, delaySeconds time.Duration) {
	vars := map[string]string{"SYSLOG_IDENTIFIER": syslogIdentifier}
	journal.Send(msg, msgType, vars) //nolint:errcheck
	time.Sleep(delaySeconds * time.Second)
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
			SendEventsToJournal(daemonName, fmt.Sprintf(
				"\n====== STACKTRACE ======\n%v\n%s\n====== /STACKTRACE ======\n",
				time.Now(),
				stackDump,
			), journal.PriDebug, 2)
		}

		os.Exit(1)
	}()
}

func DeferFuncForRunLogDriver() {
	if LoggerErr != nil {
		SendEventsToJournal(daemonName, LoggerErr.Error(), journal.PriErr, 1)
	}
}
