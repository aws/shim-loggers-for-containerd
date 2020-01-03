package logger

import (
	"sync"

	"github.com/aws/shim-loggers-for-containerd/debug"

	"github.com/coreos/go-systemd/journal"
)

// bufferedLogger is a wrapper of underlying log driver and an
// intermediate buffer between container pipes and underlying log
// driver.
type bufferedLogger struct {
	l      LogDriver
	buffer *logBuffer
}

type logBuffer struct {
	lock sync.Mutex
	wait *sync.Cond

	maxSizeInBytes int
	isClosed       bool
	queue          [][]byte
}

// NewBufferedLogger creates a logger with the provided loggerOpts
// and a buffer with customized max size.
func NewBufferedLogger(l LogDriver, maxBufferSize int) LogDriver {
	return &bufferedLogger{
		l:      l,
		buffer: newLoggerBuffer(maxBufferSize),
	}
}

// newLoggerBuffer creates a buffer that stores messages which are
// from container and consumed by sub-level log drivers.
func newLoggerBuffer(maxBufferSize int) *logBuffer {
	lb := &logBuffer{
		maxSizeInBytes: maxBufferSize,
		queue:          make([][]byte, maxBufferSize),
		isClosed:       false,
	}
	lb.wait = sync.NewCond(&lb.lock)

	return lb
}

// Start starts the non-blocking mode logger
func (bl *bufferedLogger) Start(ready func() error) error {
	// TODO: implement non-blocking mode
	debug.SendEventsToJournal(DaemonName, "Skipping non-blocking mode for now...", journal.PriInfo)
	return bl.l.Start(ready) // directly start blocking mode log driver for now
}
