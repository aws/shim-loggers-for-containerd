// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build unit
// +build unit

package logger

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	dockerlogger "github.com/docker/docker/daemon/logger"
	"github.com/stretchr/testify/require"
)

// recordingClient records every msg.Line it receives, copying the bytes
// to decouple from any upstream buffer reuse.
type recordingClient struct {
	mu    sync.Mutex
	lines []string
}

func (r *recordingClient) Log(msg *dockerlogger.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lines = append(r.lines, string(msg.Line))
	return nil
}

// slowClient simulates a slow destination, forcing the ring buffer to fill.
type slowClient struct {
	mu    sync.Mutex
	count int
	delay time.Duration
}

func (s *slowClient) Log(_ *dockerlogger.Message) error {
	if s.delay > 0 {
		time.Sleep(s.delay)
	}
	s.mu.Lock()
	s.count++
	s.mu.Unlock()
	return nil
}

// --- Binary size scenario ---

// TestBinarySize builds the shim-logger binary and asserts its size stays below
// a threshold. Update maxBinarySize downward as optimizations land.
func TestBinarySize(t *testing.T) {
	const (
		// Binary size varies by OS/arch (~26-27 MiB without ldflags).
		// Threshold set with headroom to accommodate different targets.
		// Ratchet down after applying -ldflags="-s -w" (expect ~18 MiB).
		maxBinarySize int64 = 35 * 1024 * 1024 // 35 MiB
	)

	tmpBinary, err := os.CreateTemp("", "shim-logger-size-test-*")
	require.NoError(t, err)
	tmpPath := tmpBinary.Name()
	require.NoError(t, tmpBinary.Close())
	defer func() { _ = os.Remove(tmpPath) }()

	cmd := exec.Command("go", "build", "-o", tmpPath, ".") //nolint:gosec // build path is a temp file, not user input
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %s", string(out))

	info, err := os.Stat(tmpPath)
	require.NoError(t, err)

	t.Logf("binary size = %d bytes (%.1f MiB)", info.Size(), float64(info.Size())/(1024*1024))

	if info.Size() > maxBinarySize {
		t.Errorf("binary size %d bytes (%.1f MiB) exceeds threshold %d bytes (%.1f MiB)",
			info.Size(), float64(info.Size())/(1024*1024),
			maxBinarySize, float64(maxBinarySize)/(1024*1024))
	}
}

// --- Data integrity scenario ---

// TestNoDataCorruptionThroughRingBuffer sends uniquely-identifiable messages
// through the ring buffer and verifies every received message arrives intact.
func TestNoDataCorruptionThroughRingBuffer(t *testing.T) {
	const (
		numMessages   = 10_000
		lineSize      = 200
		maxBufferSize = 1 * 1024 * 1024 // 1 MiB — small enough to cause churn
	)

	recorder := &recordingClient{}

	var pipeBuf bytes.Buffer
	expected := make([]string, numMessages)
	for i := 0; i < numMessages; i++ {
		msg := fmt.Sprintf("msg-%05d-%s", i, strings.Repeat("x", lineSize-10))
		if len(msg) > lineSize {
			msg = msg[:lineSize]
		}
		expected[i] = msg
		pipeBuf.WriteString(msg)
		pipeBuf.WriteByte('\n')
	}
	var emptyPipe bytes.Buffer

	inner, err := NewLogger(
		WithStdout(&pipeBuf),
		WithStderr(&emptyPipe),
		WithStream(recorder),
		WithInfo(NewInfo(testContainerID, testContainerName)),
	)
	require.NoError(t, err)

	bl := NewBufferedLogger(inner, DefaultBufSizeInBytes, maxBufferSize, testContainerID)

	cleanupTime := 2 * time.Second
	err = bl.Start(context.Background(), &cleanupTime, func() error { return nil })
	require.NoError(t, err)

	recorder.mu.Lock()
	received := recorder.lines
	recorder.mu.Unlock()

	t.Logf("sent %d messages, received %d", numMessages, len(received))
	require.Greater(t, len(received), 0, "expected at least some messages to be received")

	expectedSet := make(map[string]bool, numMessages)
	for _, e := range expected {
		expectedSet[e] = true
	}

	for i, line := range received {
		if !expectedSet[line] {
			t.Errorf("message %d has unexpected content (possible corruption): %q",
				i, line[:min(80, len(line))])
		}
	}
}

// --- Pointer retention scenario ---

// TestHeapAfterDrain demonstrates the pointer retention problem in the ring
// buffer. When the queue always has at least one message (normal operating
// state), the backing array retains pointers to ALL previously-dequeued
// messages, preventing GC from collecting them.
//
// The test enqueues many messages, then dequeues all but one. With the bug,
// the ~80 MiB of dequeued message data is retained. With the fix (nil-out
// on dequeue), GC reclaims it.
func TestHeapAfterDrain(t *testing.T) {
	const (
		lineSize      = 1024              // 1 KiB per message
		maxBufferSize = 100 * 1024 * 1024 // 100 MiB — large enough to accept all messages
		numMessages   = 80_000            // ~80 MiB of line data
		// After draining all but one message and GC, the retained heap should
		// be close to 1 message (~1 KiB) plus runtime overhead.
		// With the bug, it retains ~80 MiB.
		maxHeapAfterDrain uint64 = 2 * 1024 * 1024 // 2 MiB
	)

	runtime.GC()
	runtime.GC()
	var baseline runtime.MemStats
	runtime.ReadMemStats(&baseline)

	rb := newLoggerBuffer(maxBufferSize)

	for i := 0; i < numMessages; i++ {
		msg := &dockerlogger.Message{
			Line:      []byte(strings.Repeat("x", lineSize)),
			Timestamp: dummyTime,
			Source:    "stdout",
		}
		_ = rb.Enqueue(msg)
	}

	queued := len(rb.queue)
	t.Logf("buffer accepted %d of %d messages", queued, numMessages)
	require.Equal(t, numMessages, queued)

	// Drain all but one — this keeps the backing array alive (cap > 0),
	// which is the normal operating state of the ring buffer.
	for len(rb.queue) > 1 {
		_, _ = rb.Dequeue()
	}
	require.Equal(t, 1, len(rb.queue))

	runtime.GC()
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	// Ensure the ring buffer (and its backing array) is still live.
	runtime.KeepAlive(rb)

	var retained uint64
	if after.HeapInuse > baseline.HeapInuse {
		retained = after.HeapInuse - baseline.HeapInuse
	}

	t.Logf("heap baseline = %.1f MiB, heap after drain = %.1f MiB, retained = %.1f MiB",
		float64(baseline.HeapInuse)/(1024*1024),
		float64(after.HeapInuse)/(1024*1024),
		float64(retained)/(1024*1024))

	if retained > maxHeapAfterDrain {
		t.Errorf("retained heap after drain = %d bytes (%.1f MiB) exceeds threshold %d bytes (%.1f MiB) — "+
			"dequeued messages are still referenced by the backing array",
			retained, float64(retained)/(1024*1024),
			maxHeapAfterDrain, float64(maxHeapAfterDrain)/(1024*1024))
	}
}

// --- Heap memory scenarios ---

// measurePeakHeap runs a non-blocking logger workload with a slow destination
// and samples HeapInuse while the ring buffer is under pressure.
func measurePeakHeap(t *testing.T, lineSize, maxBufferSize, numMessages int, destDelay time.Duration) uint64 {
	t.Helper()

	dest := &slowClient{delay: destDelay}

	var pipeBuf bytes.Buffer
	line := strings.Repeat("a", lineSize)
	for i := 0; i < numMessages; i++ {
		pipeBuf.WriteString(line)
		pipeBuf.WriteByte('\n')
	}
	var emptyPipe bytes.Buffer

	inner, err := NewLogger(
		WithStdout(&pipeBuf),
		WithStderr(&emptyPipe),
		WithStream(dest),
		WithInfo(NewInfo(testContainerID, testContainerName)),
	)
	require.NoError(t, err)

	bl := NewBufferedLogger(inner, DefaultBufSizeInBytes, maxBufferSize, testContainerID)

	// Sample heap in background while the workload runs.
	var maxHeap atomic.Uint64
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				if m.HeapInuse > maxHeap.Load() {
					maxHeap.Store(m.HeapInuse)
				}
			}
		}
	}()

	cleanupTime := 1 * time.Second
	_ = bl.Start(context.Background(), &cleanupTime, func() error { return nil })
	close(done)

	// Final measurement after completion.
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if m.HeapInuse > maxHeap.Load() {
		maxHeap.Store(m.HeapInuse)
	}

	return maxHeap.Load()
}

// measurePeakHeapBlocking runs a blocking-mode logger workload with a slow
// destination and samples HeapInuse during the run.
func measurePeakHeapBlocking(t *testing.T, lineSize, bufferSizeInBytes, numMessages int, destDelay time.Duration) uint64 {
	t.Helper()

	dest := &slowClient{delay: destDelay}

	var pipeBuf bytes.Buffer
	line := strings.Repeat("a", lineSize)
	for i := 0; i < numMessages; i++ {
		pipeBuf.WriteString(line)
		pipeBuf.WriteByte('\n')
	}
	var emptyPipe bytes.Buffer

	inner, err := NewLogger(
		WithStdout(&pipeBuf),
		WithStderr(&emptyPipe),
		WithStream(dest),
		WithInfo(NewInfo(testContainerID, testContainerName)),
		WithBufferSizeInBytes(bufferSizeInBytes),
	)
	require.NoError(t, err)

	var maxHeap atomic.Uint64
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				if m.HeapInuse > maxHeap.Load() {
					maxHeap.Store(m.HeapInuse)
				}
			}
		}
	}()

	cleanupTime := 1 * time.Second
	_ = inner.Start(context.Background(), &cleanupTime, func() error { return nil })
	close(done)

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if m.HeapInuse > maxHeap.Load() {
		maxHeap.Store(m.HeapInuse)
	}

	return maxHeap.Load()
}

func TestMemoryScenario_SmallLines_Blocking_DefaultBuffer(t *testing.T) {
	const (
		lineSize    = 100
		bufferSize  = DefaultBufSizeInBytes // 16 KiB
		numMessages = 10_000
		// 10K messages × 100 bytes = ~1 MiB of data. With GC overhead and
		// runtime baseline, peak heap should stay well under 10 MiB.
		maxHeapBytes uint64 = 10 * 1024 * 1024
	)

	heap := measurePeakHeapBlocking(t, lineSize, bufferSize, numMessages, 50*time.Microsecond)
	t.Logf("SmallLines/Blocking: max HeapInuse = %d bytes (%.1f MiB)", heap, float64(heap)/(1024*1024))

	if heap > maxHeapBytes {
		t.Errorf("HeapInuse %d bytes (%.1f MiB) exceeds threshold %d bytes (%.1f MiB)",
			heap, float64(heap)/(1024*1024),
			maxHeapBytes, float64(maxHeapBytes)/(1024*1024))
	}
}

func TestMemoryScenario_SmallLines_Blocking_LargeBuffer(t *testing.T) {
	const (
		lineSize    = 100
		bufferSize  = 256 * 1024 // 256 KiB (awslogs-style)
		numMessages = 10_000
		// Same data volume as above, larger read buffer adds ~256 KiB × 2 pipes.
		maxHeapBytes uint64 = 15 * 1024 * 1024
	)

	heap := measurePeakHeapBlocking(t, lineSize, bufferSize, numMessages, 50*time.Microsecond)
	t.Logf("SmallLines/Blocking/LargeBuffer: max HeapInuse = %d bytes (%.1f MiB)", heap, float64(heap)/(1024*1024))

	if heap > maxHeapBytes {
		t.Errorf("HeapInuse %d bytes (%.1f MiB) exceeds threshold %d bytes (%.1f MiB)",
			heap, float64(heap)/(1024*1024),
			maxHeapBytes, float64(maxHeapBytes)/(1024*1024))
	}
}

func TestMemoryScenario_SmallLines_NonBlocking_LargeBuffer(t *testing.T) {
	const (
		lineSize      = 100
		maxBufferSize = 256 * 1024 * 1024 // 256 MiB
		numMessages   = 10_000
		// 10K messages × 100 bytes = ~1 MiB of data. Even with a 256 MiB buffer
		// configured, the actual data is small. Peak heap should stay under 15 MiB.
		maxHeapBytes uint64 = 15 * 1024 * 1024
	)

	heap := measurePeakHeap(t, lineSize, maxBufferSize, numMessages, 50*time.Microsecond)
	t.Logf("SmallLines/NonBlocking/LargeBuffer: max HeapInuse = %d bytes (%.1f MiB)", heap, float64(heap)/(1024*1024))

	if heap > maxHeapBytes {
		t.Errorf("HeapInuse %d bytes (%.1f MiB) exceeds threshold %d bytes (%.1f MiB)",
			heap, float64(heap)/(1024*1024),
			maxHeapBytes, float64(maxHeapBytes)/(1024*1024))
	}
}

func TestMemoryScenario_SmallLines_NonBlocking_DefaultBuffer(t *testing.T) {
	const (
		lineSize      = 100
		maxBufferSize = 10 * 1024 * 1024 // 10 MiB
		numMessages   = 10_000
		// Same as above — data volume is ~1 MiB regardless of buffer config.
		maxHeapBytes uint64 = 15 * 1024 * 1024
	)

	heap := measurePeakHeap(t, lineSize, maxBufferSize, numMessages, 50*time.Microsecond)
	t.Logf("SmallLines: max HeapInuse = %d bytes (%.1f MiB)", heap, float64(heap)/(1024*1024))

	if heap > maxHeapBytes {
		t.Errorf("HeapInuse %d bytes (%.1f MiB) exceeds threshold %d bytes (%.1f MiB)",
			heap, float64(heap)/(1024*1024),
			maxHeapBytes, float64(maxHeapBytes)/(1024*1024))
	}
}

func TestMemoryScenario_LargeLines_NonBlocking_DefaultBuffer(t *testing.T) {
	const (
		lineSize      = 62 * 1024        // 62 KiB
		maxBufferSize = 10 * 1024 * 1024 // 10 MiB
		numMessages   = 1_000
		// 1K messages × 62 KiB = ~62 MiB of data into a 10 MiB buffer.
		// Peak heap reaches ~160 MiB (16x buffer) due to GC headroom and
		// message struct overhead. With -race instrumentation, add ~30% headroom.
		maxHeapBytes uint64 = 210 * 1024 * 1024
	)

	heap := measurePeakHeap(t, lineSize, maxBufferSize, numMessages, 100*time.Microsecond)
	t.Logf("LargeLines: max HeapInuse = %d bytes (%.1f MiB)", heap, float64(heap)/(1024*1024))

	if heap > maxHeapBytes {
		t.Errorf("HeapInuse %d bytes (%.1f MiB) exceeds threshold %d bytes (%.1f MiB)",
			heap, float64(heap)/(1024*1024),
			maxHeapBytes, float64(maxHeapBytes)/(1024*1024))
	}
}
