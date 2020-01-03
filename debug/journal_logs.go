package debug

import (
	"github.com/coreos/go-systemd/journal"
)

// This is a temporary solution for logging the shim-loggers-for-containerd package itself. We directly
// send the events to system journal and they are identified by the package name. Since this process is
// started by containerd, we can check the logs using `journalctl -u containerd.service`.
func SendEventsToJournal(syslogIdentifier string, msg string, msgType journal.Priority) {
	vars := map[string]string{"SYSLOG_IDENTIFIER": syslogIdentifier}
	journal.Send(msg, msgType, vars)
}
