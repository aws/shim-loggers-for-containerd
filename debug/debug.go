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
	Verbose = false
)

// This is a temporary solution for logging the shim-loggers-for-containerd package itself. We directly
// send the events to system journal and they are identified by the package name. Since this process is
// started by containerd, we can check the logs using `journalctl -u containerd.service`.
func SendEventsToJournal(syslogIdentifier string, msg string, msgType journal.Priority) {
	vars := map[string]string{"SYSLOG_IDENTIFIER": syslogIdentifier}
	journal.Send(msg, msgType, vars)
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
			), journal.PriDebug)
		}

		os.Exit(1)
	}()
}
