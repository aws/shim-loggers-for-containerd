package logger

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/aws/shim-loggers-for-containerd/debug"

	"github.com/coreos/go-systemd/journal"
	"github.com/pkg/errors"
)

// bufferedLogger is a wrapper of underlying log driver and an
// intermediate buffer between container pipes and underlying log
// driver.
type bufferedLogger struct {
	l      LogDriver
	buffer *logBuffer
	// closedPipe is a channel listens to goroutine who sends logs from buffer
	// to destination to get notification from the other two goroutines, who scan
	// the container io pipes, that both the pipes are closed. At that time,
	// logger flushes all the messages left in the buffer to destination and
	// clear the buffer.
	closedPipes chan bool
}

// Adopted from https://github.com/moby/moby/blob/master/daemon/logger/ring.go#L128
// as this struct is not exported.
type logBuffer struct {
	// A mutex lock is used here when writing/reading log messages from the queue
	// as there exists three go routines accessing the buffer.
	lock sync.Mutex
	// A condition variable wait is used here to notify goroutines that get access to
	// the buffer should wait or continue.
	wait *sync.Cond
	// current total bytes stored in the buffer
	curSizeInBytes int
	// maximum bytes capacity provided by the buffer
	maxSizeInBytes int
	// queue saves all the log messages read from pipes exposed by containerd, and
	// is consumed by underlying log driver.
	queue []*msg
}

// msg stores a single line of log message, source pipe name (stdout/stderr),
// and the timestamp that the log obtained from the pipe.
type msg struct {
	line    []byte
	source  string
	logTime time.Time
}

// NewBufferedLogger creates a logger with the provided LoggerOpt,
// a buffer with customized max size and a channel monitor if stdout
// and stderr pipes are closed.
func NewBufferedLogger(l LogDriver, maxBufferSize int) LogDriver {
	return &bufferedLogger{
		l:           l,
		buffer:      newLoggerBuffer(maxBufferSize),
		closedPipes: make(chan bool, 2),
	}
}

// newLoggerBuffer creates a buffer that stores messages which are
// from container and consumed by sub-level log drivers.
func newLoggerBuffer(maxBufferSize int) *logBuffer {
	lb := &logBuffer{
		maxSizeInBytes: maxBufferSize,
		queue:          make([]*msg, 0),
	}
	lb.wait = sync.NewCond(&lb.lock)

	return lb
}

// Start starts the non-blocking mode logger.
func (bl *bufferedLogger) Start(uid int, gid int, ready func() error) error {
	var wg sync.WaitGroup
	stdout, stderr := bl.l.GetPipes()
	if stdout != nil {
		wg.Add(1)
		debug.SendEventsToJournal(DaemonName,
			"Starting reading from stdout pipe",
			journal.PriInfo)
		go bl.saveLogsToBuffer(stdout, &wg, sourceSTDOUT, uid, gid)
	}
	if stderr != nil {
		wg.Add(1)
		debug.SendEventsToJournal(DaemonName,
			"Starting reading from stderr pipe",
			journal.PriInfo)
		go bl.saveLogsToBuffer(stderr, &wg, sourceSTDERR, uid, gid)
	}

	// Signal that the container is ready to be started
	if err := ready(); err != nil {
		debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr)
		return errors.Wrap(err, "failed to notify container ready status to containerd")
	}

	// Start the underling log driver to send logs to destination
	wg.Add(1)
	go bl.sendLogs(&wg, uid, gid)

	wg.Wait()
	debug.SendEventsToJournal(DaemonName,
		"All logs saved to buffer has been sent to destination",
		journal.PriInfo)

	return nil
}

// saveLogsToBuffer saves container logs to intermediate buffer.
func (bl *bufferedLogger) saveLogsToBuffer(f io.Reader, wg *sync.WaitGroup, source string, uid int, gid int) {
	defer wg.Done()

	// Set uid for this goroutine. Currently the Setuid syscall does not
	// apply on threads in golang, see issue: https://github.com/golang/go/issues/1435
	// TODO: remove it once the changes are released: https://go-review.googlesource.com/c/go/+/210639
	if err := SetUIDAndGID(uid, gid); err != nil {
		debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr)
		return
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	// Scan function breaks either when catch EOF error or there exists error
	// scanning message except for EOF error.
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			debug.SendEventsToJournal(DaemonName,
				"Message is empty, skip saving", journal.PriInfo)
			continue
		}
		if err := bl.read(scanner, source); err != nil {
			debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr)
			return
		}
	}

	// No messages in the pipe, send signal to closed pipe channel
	debug.SendEventsToJournal(DaemonName, "Pipe closed", journal.PriInfo)
	bl.closedPipes <- true
}

// read reads a single log message from pipe and saves it to buffer.
func (bl *bufferedLogger) read(s *bufio.Scanner, source string) error {
	if s.Err() != nil {
		return errors.Wrap(s.Err(), "failed to get logs from container")
	}
	if debug.Verbose {
		debug.SendEventsToJournal(DaemonName,
			fmt.Sprintf("[SCANNER] Scanned msg: %s", s.Text()),
			journal.PriDebug)
		debug.SendEventsToJournal(DaemonName,
			fmt.Sprintf("current buffer size: %d", bl.buffer.curSizeInBytes),
			journal.PriDebug)
	}

	logMsg := &msg{
		line:    s.Bytes(),
		source:  source,
		logTime: time.Now(),
	}
	err := bl.buffer.Enqueue(logMsg)
	if err != nil {
		return errors.Wrap(s.Err(), "failed to save logs to buffer")
	}

	return nil
}

// sendLogs consumes logs from intermediate buffer and use the
// underlying log drive to send logs to destination.
func (bl *bufferedLogger) sendLogs(wg *sync.WaitGroup, uid int, gid int) {
	defer wg.Done()

	// Set uid for this goroutine. Currently the Setuid syscall does not
	// apply on threads in golang, see issue: https://github.com/golang/go/issues/1435
	// TODO: remove it once the changes are released: https://go-review.googlesource.com/c/go/+/210639
	if err := SetUIDAndGID(uid, gid); err != nil {
		debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr)
		return
	}

	count := 0
	for {
		select {
		case <-bl.closedPipes:
			count++
			// If received closed pipe signal from both pipes,
			// flush messages left in buffer.
			if count == 2 {
				debug.SendEventsToJournal(DaemonName,
					"All pipes are closed, closing buffer",
					journal.PriInfo)
				if err := bl.flushMessages(); err != nil {
					debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr)
					return
				}
				return
			}
		default:
			if err := bl.send(); err != nil {
				debug.SendEventsToJournal(DaemonName, err.Error(), journal.PriErr)
				return
			}
		}
	}
}

// send dequeues a single log message from buffer and sends to destination.
func (bl *bufferedLogger) send() error {
	msg, err := bl.buffer.Dequeue()
	if err != nil {
		return errors.Wrap(err, "failed to read logs from buffer")
	}

	err = bl.LogWithRetry(msg.line, msg.source, msg.logTime)
	if err != nil {
		return errors.Wrap(err, "failed to send logs to destination")
	}

	return nil
}

// flushMessages flushes all the messages left in the buffer to destination
// after container pipes are closed.
func (bl *bufferedLogger) flushMessages() error {
	messages := bl.buffer.Flush()
	for _, msg := range messages {
		err := bl.LogWithRetry(msg.line, msg.source, msg.logTime)
		if err != nil {
			return errors.Wrap(err, "unable to flush the remaining messages to destination")
		}
	}

	return nil
}

// GetPipes gets pipes of container that exposed by containerd.
func (bl *bufferedLogger) GetPipes() (io.Reader, io.Reader) {
	return bl.l.GetPipes()
}

// LogWithRetry lets underlying log driver send logs to destination.
func (bl *bufferedLogger) LogWithRetry(line []byte, source string, logTimestamp time.Time) error {
	if debug.Verbose {
		debug.SendEventsToJournal(DaemonName,
			fmt.Sprintf("[BUFFER] Sending message: %s", string(line)),
			journal.PriDebug)
	}
	return bl.l.LogWithRetry(line, source, logTimestamp)
}

// Adopted from https://github.com/moby/moby/blob/master/daemon/logger/ring.go#L155
// as messageRing struct is not exported.
// Enqueue adds a single log message to the tail of intermediate buffer.
func (b *logBuffer) Enqueue(msg *msg) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	lineSizeInBytes := len(msg.line)
	// If there is not enough space left for the new coming log message or the message
	// size is larger than the whole buffer size.
	if b.curSizeInBytes+lineSizeInBytes > b.maxSizeInBytes ||
		lineSizeInBytes > b.maxSizeInBytes {
		if debug.Verbose {
			debug.SendEventsToJournal(DaemonName,
				"buffer is full/message is too long, waiting for available bytes",
				journal.PriDebug)
			debug.SendEventsToJournal(DaemonName,
				fmt.Sprintf("message size: %d, current buffer size: %d, max buffer size %d",
					lineSizeInBytes,
					b.curSizeInBytes,
					b.maxSizeInBytes),
				journal.PriDebug)
		}

		// Wake up "Dequeue" or the other "Enqueue" go routine (called by the other pipe)
		// waiting on current mutex lock if there's any
		b.wait.Signal()
		return nil
	}

	b.queue = append(b.queue, msg)
	b.curSizeInBytes += lineSizeInBytes
	// Wake up "Dequeue" or the other "Enqueue" go routine (called by the other pipe)
	// waiting on current mutex lock if there's any
	b.wait.Signal()

	return nil
}

// Adopted from https://github.com/moby/moby/blob/master/daemon/logger/ring.go#L179
// as messageRing struct is not exported.
// Dequeue gets a line of log message from the head of intermediate buffer.
func (b *logBuffer) Dequeue() (*msg, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// If there is no logs yet in the buffer, wait suspends current go routine
	for len(b.queue) == 0 {
		if debug.Verbose {
			debug.SendEventsToJournal(DaemonName,
				"No messages in queue, waiting...",
				journal.PriDebug)
		}
		b.wait.Wait()
	}

	// Get and remove the oldest message saved in buffer/queue from head and update
	// the current used bytes of buffer.
	msg := b.queue[0]
	b.queue = b.queue[1:]
	b.curSizeInBytes -= len(msg.line)

	return msg, nil
}

// Adopted from https://github.com/moby/moby/blob/master/daemon/logger/ring.go#L215
// as messageRing struct is not exported.
// Flush flushes all the messages left in the buffer and clear queue.
func (b *logBuffer) Flush() []*msg {
	b.lock.Lock()
	defer b.lock.Unlock()

	if len(b.queue) == 0 {
		return make([]*msg, 0)
	}

	messages := b.queue
	b.queue = make([]*msg, 0)

	return messages
}
