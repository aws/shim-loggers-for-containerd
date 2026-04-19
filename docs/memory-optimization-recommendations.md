<!-- markdownlint-disable MD024 -->
# Shim-Logger Memory Optimization Recommendations

## Context

The shim-logger forwards container stdout/stderr to destinations such as CloudWatch Logs,
Splunk, and Fluentd. Each container with logging configured runs its own shim-logger process.
Benchmarking shows RSS memory can reach 1.6+ GiB in extreme cases (512 MiB non-blocking buffer,
62 KiB log lines), and in production has been observed as high as 4 GiB.

As of June 25, 2025, ECS changed the default log driver mode from `blocking` to `non-blocking`
with a default `max-buffer-size` of `10m`. This means every Fargate task with logging configured
now runs the ring buffer code path by default, making these optimizations relevant to all
customers rather than just those who explicitly opted into non-blocking mode.

Reference: [ECS LogConfiguration API](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_LogConfiguration.html)

---

## Recommendation 1: Strip Debug Symbols from the Binary ✅ Implemented

### Impact

Reduces the compiled binary size by approximately 20-30%. Because the Go runtime memory-maps the
binary into the process address space, a smaller binary directly reduces the baseline RSS
footprint before any log data is processed. This affects every single shim-logger instance.

### Why This Works

The Go toolchain embeds two categories of metadata into compiled binaries by default:

1. **Symbol table** (`-s` flag removes it): The symbol table maps memory addresses back to
   function and variable names. It is used by debuggers (e.g., `dlv`, `gdb`) and profiling tools
   to translate raw addresses into human-readable stack traces. In a production deployment where
   the shim-logger is not being actively debugged, this table serves no runtime purpose but still
   occupies space in the binary's ELF/Mach-O sections, which get mapped into the process's
   virtual address space on startup.

2. **DWARF debugging information** (`-w` flag removes it): DWARF is a standardized debugging
   data format that encodes detailed type information, variable locations, line number mappings,
   and inlining records. It enables debuggers to step through source code, inspect variables, and
   reconstruct call stacks. This data can be substantial — often larger than the symbol table
   itself — because it encodes the full type graph of the program and all its dependencies.

These flags are passed to the Go linker via `-ldflags`. The current `Makefile` build target is:

```makefile
$(AWS_CONTAINERD_LOGGERS_BINARY):
    go build -o $(AWS_CONTAINERD_LOGGERS_BINARY) $(AWS_CONTAINERD_LOGGERS_DIR)
```

The recommended change is:

```makefile
$(AWS_CONTAINERD_LOGGERS_BINARY):
    go build -ldflags="-s -w" -o $(AWS_CONTAINERD_LOGGERS_BINARY) $(AWS_CONTAINERD_LOGGERS_DIR)
```

The production build process calls `make build`, which delegates directly to this Makefile
target. No additional linker flags are injected. Updating the Makefile is sufficient.

When the OS loads the binary, it maps the text (code), rodata (read-only data), and other
sections into the process's virtual memory. The symbol table and DWARF sections are part of this
mapping. Removing them means fewer pages are mapped, which directly reduces the resident set size
reported by the kernel. The Go runtime itself does not use these sections at runtime — they exist
solely for external tooling.

### Trade-offs

- **Loss of debuggability**: If a production crash occurs, stack traces from `panic` output will
  still contain function names (Go embeds a separate, minimal `gopclntab` for this purpose), but
  external debuggers will not be able to map addresses to source lines or inspect variable values.
  If post-mortem debugging with `dlv` or `gdb` is part of the operational workflow, this change
  would impair that. A mitigation is to keep unstripped binaries in the build artifact store and
  only deploy stripped binaries.
- **No functional impact**: The `-s -w` flags do not alter code generation, optimization levels,
  or runtime behavior. The compiled machine code is identical; only metadata sections are removed.
- **No customer impact**: Customers do not interact with the shim-logger binary directly. They
  configure logging via the ECS task definition. Stripping symbols is invisible to them.

---

## Recommendation 2: Nil Out Dequeued Pointers ✅ Implemented

### Impact

Reduces steady-state heap size in non-blocking mode by allowing the garbage collector to reclaim
`*dockerlogger.Message` structs as soon as they are consumed, rather than retaining them until
the queue slice's backing array is reallocated. For workloads with high message throughput, this
can reduce heap retention by the size of all previously-dequeued-but-still-referenced messages.

### Why This Works

The ring buffer's `Dequeue` method in `logger/buffered_logger.go` currently removes messages
from the head of the queue by reslicing:

```go
msg := b.queue[0]
b.queue = b.queue[1:]
b.curSizeInBytes -= len(msg.Line)
```

This reslice moves the start pointer of the slice forward, but the underlying backing array
still holds a pointer to the dequeued `*dockerlogger.Message` at its original index. The Go
garbage collector cannot collect an object if any live reference points to it. Because the
backing array is a single contiguous allocation, all pointers within it — including those
"before" the current slice start — remain reachable as long as the slice header exists.

In practice, this means that if 10,000 messages have been enqueued and 9,000 have been dequeued,
the backing array still holds pointers to all 10,000 messages. The 9,000 dequeued messages and
their `Line` byte slices cannot be garbage collected.

The fix is to nil out the pointer before reslicing:

```go
msg := b.queue[0]
b.queue[0] = nil        // allow GC to collect the message
b.queue = b.queue[1:]
b.curSizeInBytes -= len(msg.Line)
```

Additionally, Docker's logger package provides a `sync.Pool`-based message recycling mechanism
via `dockerlogger.NewMessage()` (used in `newMessage` in `common.go`) and a corresponding
`PutMessage()`. An audit of the Docker driver source (see Trade-offs below) reveals that all
three drivers (awslogs, splunk, fluentd) already call `PutMessage` internally after consuming
the message. The shim-logger must NOT call `PutMessage` itself, as that would double-return the
message to the pool and risk data corruption. The shim-logger's optimization is limited to the
nil-out fix shown above.

### Trade-offs

- **Nil-out is zero-risk**: A single pointer write with no functional side effects.
- **Do NOT call `PutMessage` from the shim-logger**: All three Docker drivers (awslogs, splunk,
  fluentd) already call `logger.PutMessage(msg)` internally after copying `msg.Line`. Calling it
  again from the shim-logger would double-return the message to the pool, causing data
  corruption. The shim-logger's optimization is limited to the nil-out fix in `Dequeue`.
- **No customer impact**: Log delivery behavior, ordering, and loss characteristics are unchanged.

### Safeguard: Data Integrity Test

A test must verify that log content is never corrupted by pool reuse. This catches double
`PutMessage`, premature `PutMessage`, or any future change to the Docker drivers' pool behavior.

Location: `logger/pool_lifecycle_test.go` with a `//go:build unit` tag.

```go
func TestNoDataCorruptionWithPoolReuse(t *testing.T) {
    // 1. Create a mock stream that records every msg.Line it receives.
    // 2. Create a bufferedLogger with the mock stream.
    // 3. Send 10,000 messages with unique content (e.g., "msg-00001").
    // 4. After all messages are processed, verify every recorded Line matches
    //    its expected content exactly.
    // 5. If any Line contains content from a different message, fail.
}
```

This test should run on every CI build. If a Docker version upgrade changes the `PutMessage`
contract, this test will catch the resulting corruption.

---

## Future Considerations

The following recommendations have not been implemented. They involve runtime tuning decisions
or behavioral changes that should be evaluated against production workload data after the
impact of Recommendations 1 and 2 is measured in production.

---

## Recommendation 3: Set GOMEMLIMIT Based on Configured Buffer Size

### Impact

Caps the maximum Go heap size, preventing the garbage collector from allowing the heap to grow
unboundedly. This directly limits worst-case RSS and closes the gap between the configured
`max-buffer-size` and the actual observed memory usage.

### Why This Works

By default, the Go garbage collector uses the `GOGC` environment variable (default: 100) to
determine when to trigger a collection. `GOGC=100` means the GC triggers when the heap has grown
to 2x the size of the live set from the previous collection. This means if the live set is
50 MiB (ring buffer contents + message structs + read buffers), the heap can grow to 100 MiB
before GC runs. After collection, if the live set is now 60 MiB, the next GC target is 120 MiB.
This feedback loop allows the heap to grow significantly beyond the actual data being held.

Go 1.19 introduced `GOMEMLIMIT`, which sets a soft memory limit for the Go heap. When set, the
GC will run more aggressively as the heap approaches this limit, regardless of the `GOGC`
setting. This is specifically designed for workloads where memory is a constrained resource —
exactly the shim-logger's situation.

The shim-logger knows its memory budget at startup:

- The `max-buffer-size` flag defines the ring buffer capacity (e.g., 10 MiB default).
- The read buffers are fixed: 2 pipes × `bufferSizeInBytes` (256 KiB for awslogs, 16 KiB for
  others).
- A baseline overhead for the Go runtime, goroutine stacks, and the binary's mapped sections.

A reasonable formula:

```go
import "runtime/debug"

baselineBytes := 20 * 1024 * 1024  // 20 MiB for runtime, stacks, binary overhead
readBufferBytes := 2 * bufferSizeInBytes
totalLimit := int64(baselineBytes + maxBufferSize + readBufferBytes)

// Add 50% headroom for transient allocations (message structs in flight,
// SDK request/response buffers, TLS buffers, etc.)
debug.SetMemoryLimit(totalLimit * 3 / 2)
```

This should be called early in `main()` after parsing flags but before starting the logger.

### Trade-offs

- **Increased GC CPU usage**: When the heap approaches `GOMEMLIMIT`, the GC runs more
  frequently. For a shim-logger that is I/O-bound (waiting on pipe reads and network writes to
  CloudWatch/Splunk/Fluentd), the additional CPU cost is generally negligible. However, under
  extreme throughput (container emitting data as fast as possible), the GC could become a
  bottleneck. The 50% headroom in the formula above mitigates this by giving the GC room to
  breathe.
- **GC thrashing risk**: If the limit is set too low, the GC will run continuously without being
  able to free enough memory, consuming 100% of available CPU. This is called "GC thrashing" and
  is the primary risk of `GOMEMLIMIT`. The formula above avoids this by sizing the limit to the
  known memory budget plus headroom. The Go runtime also has a built-in safeguard: if GC is
  consuming more than 50% of CPU time, it backs off regardless of the memory limit.
- **Soft limit, not hard limit**: `GOMEMLIMIT` is a soft limit. The Go runtime can exceed it
  temporarily if a single allocation requires it (e.g., a large SDK response buffer). It will
  not OOM-kill the process. The kernel's cgroup memory limit remains the hard boundary.
- **Customer impact**: None. This is an internal runtime tuning parameter. Log delivery behavior
  is unchanged. The only observable effect is lower RSS, which benefits customers by leaving more
  memory available for their application containers within the task's memory limit.

---

## Recommendation 4: Account for Per-Message Overhead in Ring Buffer Size Tracking

### Impact

Makes the ring buffer's memory accounting accurate, preventing actual memory usage from
exceeding the configured `max-buffer-size` by a significant margin. With the current 10 MiB
default buffer and small log lines, the untracked overhead can add 50-100% on top of the
configured size.

### Why This Works

The `Enqueue` method in `logger/buffered_logger.go` tracks buffer usage using only the line
content:

```go
lineSizeInBytes := len(msg.Line)
// ...
b.curSizeInBytes += lineSizeInBytes
```

But each `*dockerlogger.Message` struct carries additional memory beyond `msg.Line`:

| Field | Size (64-bit) | Notes |
|-------|---------------|-------|
| `Line` (slice header) | 24 bytes | pointer + length + capacity |
| `Source` (string) | 16 bytes | pointer + length |
| `Timestamp` (time.Time) | 24 bytes | wall + ext + loc pointer |
| `Attrs` (slice header) | 24 bytes | even if nil/empty |
| `PLogMetaData` (pointer) | 8 bytes | nil for non-partial messages |
| Struct padding/alignment | ~8 bytes | |
| **Total struct overhead** | **~104 bytes** | **per message, excluding Line content** |

For partial messages, `PLogMetaData` points to an additional struct:

| Field | Size | Notes |
|-------|------|-------|
| `ID` (string) | 16 bytes | header; content is 64-byte hex string = 80 bytes total |
| `Ordinal` (int) | 8 bytes | |
| `Last` (bool) | 1 byte + 7 padding | |
| **Total** | **~96 bytes** | **additional per partial message** |

Additionally, `msg.Line` is created via `append(msg.Line, line...)` in `newMessage`, which
allocates a new backing array. The Go runtime's memory allocator rounds up allocations to size
classes (8, 16, 32, 48, 64, 80, 96, ..., 32768 bytes, then page-aligned). A 100-byte line
actually consumes 112 bytes (the next size class). This rounding is not tracked.

**Concrete example**: With 10 MiB buffer and 100-byte log lines:

- Messages that fit: `10 * 1024 * 1024 / 100 = 104,857` messages (by current accounting)
- Actual memory per message: 100 (line) + 104 (struct) + 12 (allocator rounding) = ~216 bytes
- Actual memory used: `104,857 * 216 = ~21.6 MiB` — more than 2x the configured 10 MiB

The fix is to include a per-message overhead constant in the size tracking:

```go
const perMessageOverheadBytes = 104  // struct fields excluding Line content

lineSizeInBytes := len(msg.Line)
messageSizeInBytes := lineSizeInBytes + perMessageOverheadBytes

if len(b.queue) > 0 && b.curSizeInBytes+messageSizeInBytes > b.maxSizeInBytes {
    // buffer full, drop message
    b.wait.Signal()
    return nil
}

b.queue = append(b.queue, msg)
b.curSizeInBytes += messageSizeInBytes
```

And correspondingly in `Dequeue`:

```go
b.curSizeInBytes -= (len(msg.Line) + perMessageOverheadBytes)
```

### Trade-offs

- **Reduced effective buffer capacity for small messages**: With the overhead accounted for, a
  10 MiB buffer will hold fewer messages than before. For 100-byte lines, the buffer would hold
  ~48,000 messages instead of ~104,000. This means messages will be dropped sooner under
  sustained load. However, this is the *correct* behavior — the previous accounting was
  under-reporting memory usage, which meant the process was using more memory than the customer
  configured.
- **No impact for large messages**: For 62 KiB log lines, the 104-byte overhead is 0.16% of the
  message size — negligible. This optimization primarily affects workloads with many small log
  lines.
- **Customer-visible behavior change**: Customers with small, high-frequency log lines may
  observe more log drops in non-blocking mode than before, because the buffer fills up sooner
  (in terms of tracked bytes). This is a correctness improvement — the buffer was already
  consuming that memory, it just wasn't being tracked. The alternative is that the process
  exceeds its memory budget and gets OOM-killed, which is worse than dropping logs.
- **Constant may drift**: The 104-byte overhead is based on the current `dockerlogger.Message`
  struct layout. If Docker changes the struct in a future version, the constant would need
  updating. A `unsafe.Sizeof` call could be used instead, but that only captures the struct
  size, not the backing array overhead.

---

## Recommendation 5: Use Dynamic Read Buffer Sizing for awslogs

### Impact

Reduces per-pipe memory allocation from 256 KiB to 16 KiB for the common case (log lines under
16 KiB), saving ~480 KiB per shim-logger instance (2 pipes × 240 KiB savings). This is modest
in absolute terms but meaningful when multiplied across many containers on a host.

### Why This Works

The awslogs driver sets its read buffer to `maximumBytesPerEvent` (262,118 bytes ≈ 256 KiB) in
`logger/awslogs/logger.go`:

```go
logger.WithBufferSizeInBytes(maximumBytesPerEvent),
```

This value is passed to `Read()` in `logger/common.go`, which allocates a buffer of that size:

```go
buf := make([]byte, bufferSizeInBytes)
```

The 256 KiB size exists because CloudWatch Logs has a maximum event size of 256 KiB, and the
read buffer must be large enough to hold a single complete event before sending it. However, the
vast majority of log lines are far smaller — typical application logs are 100 bytes to 4 KiB.

The `Read` function already handles the case where a log line exceeds the buffer size: it splits
the line into partial messages (the `isPartialMsg` / `partialOrdinal` logic). This means a
smaller buffer does not cause data loss — it causes large lines to be split into partials, which
are reassembled by the CloudWatch Logs agent or viewer.

Two approaches:

**Approach A: Reduce the default, accept more partials for large lines.**

Change `WithBufferSizeInBytes(maximumBytesPerEvent)` to `WithBufferSizeInBytes(DefaultBufSizeInBytes)`
(16 KiB), matching the splunk and fluentd drivers. Log lines over 16 KiB would be split into
partials. This is the simplest change.

**Approach B: Start small, grow on demand.**

Replace the fixed `[]byte` buffer with a `bytes.Buffer` or a dynamically-grown slice that starts
at 16 KiB and grows (via `append` or explicit reallocation) only when a partial message is
detected and the line hasn't ended. This preserves the ability to hold a full 256 KiB event in a
single message while only paying the memory cost when needed.

### Trade-offs

- **Approach A — More partial messages for large log lines**: Customers who emit log lines
  larger than 16 KiB (e.g., JSON payloads, stack traces) would see those lines split into
  multiple CloudWatch log events with partial metadata. This changes the appearance of logs in
  CloudWatch Logs Insights queries and the console. For most customers (small log lines), there
  is no observable change.
- **Approach B — Implementation complexity**: Dynamic buffer growth adds complexity to the `Read`
  function, which is already non-trivial with its partial message state machine. The growth
  logic must handle the case where the buffer grows mid-line and the existing bytes need to be
  preserved. This is doable but increases the surface area for bugs.
- **Neither approach affects log completeness**: All bytes from the container's stdout/stderr
  are still forwarded to CloudWatch. The difference is whether a large line arrives as one event
  or multiple partial events.

---

## Recommendation 6: Tune Go GC via GOGC

### Impact

Reduces peak RSS by triggering garbage collection more frequently, at the cost of additional CPU
time spent on GC. Setting `GOGC=50` (instead of the default 100) means the GC triggers when the
heap has grown to 1.5x the live set, rather than 2x.

### Why This Works

The `GOGC` environment variable controls the garbage collector's aggressiveness. The value
represents the percentage of heap growth relative to the live set that triggers a new GC cycle:

- `GOGC=100` (default): GC triggers when heap reaches 2x the live set (100% growth).
- `GOGC=50`: GC triggers when heap reaches 1.5x the live set (50% growth).
- `GOGC=25`: GC triggers at 1.25x the live set.

For the shim-logger, the live set is dominated by the ring buffer contents (in non-blocking
mode) or the in-flight message (in blocking mode), plus the read buffers and runtime overhead.
The "garbage" is primarily dequeued messages waiting to be collected and transient allocations
from the AWS SDK (HTTP request/response buffers, TLS state, JSON serialization buffers).

With `GOGC=100`, if the live set is 10 MiB, the heap can grow to 20 MiB before GC runs. The
10 MiB of garbage represents dequeued messages and SDK temporaries that are no longer needed but
haven't been collected yet. With `GOGC=50`, the heap only grows to 15 MiB before GC runs,
reclaiming that garbage sooner.

This can be set programmatically:

```go
import "runtime/debug"

debug.SetGCPercent(50)
```

Or via environment variable: `GOGC=50`.

Note: If Recommendation 3 (`GOMEMLIMIT`) is implemented, `GOGC` tuning becomes less critical
because `GOMEMLIMIT` provides a hard ceiling that overrides the `GOGC`-based trigger when the
heap approaches the limit. However, `GOGC` still controls GC behavior when the heap is well
below the limit, so both can be used together.

### Trade-offs

- **Increased CPU usage**: More frequent GC cycles mean more CPU time spent scanning and
  collecting. For the shim-logger, which spends most of its time blocked on I/O (reading from
  pipes, writing to network), the CPU overhead is typically small. However, under extreme
  throughput (container emitting data as fast as possible with the destination unreachable,
  causing rapid buffer churn), the additional GC work could become noticeable.
- **Diminishing returns with GOMEMLIMIT**: If `GOMEMLIMIT` is set (Recommendation 3), the GC
  already runs aggressively near the limit. Lowering `GOGC` on top of that provides marginal
  benefit for the near-limit case but still helps when the heap is well below the limit.
- **Not configurable by customers**: `GOGC` is set by the shim-logger process, not by the
  customer's task definition. Customers cannot override it. If a future use case requires
  different GC tuning, it would need a code change or a new flag.

---

## Recommendation 7: Reduce Dependency Bloat

### Impact

Reduces binary size and baseline RSS by removing transitive dependencies that are linked into
the binary but never executed at runtime. This is the highest-effort recommendation but has the
largest long-term payoff for baseline memory.

### Why This Works

The Go linker includes all reachable code in the final binary. "Reachable" means any function
that is transitively referenced from `main`, even if it is never called at runtime. The
shim-logger imports three Docker logger implementations:

```go
import (
    dockerawslogs "github.com/docker/docker/daemon/logger/awslogs"
    dockerfluentd "github.com/docker/docker/daemon/logger/fluentd"
    dockersplunk  "github.com/docker/docker/daemon/logger/splunk"
)
```

These imports pull in the Docker engine's logger framework, which transitively depends on:

| Dependency | Why it's pulled in | Approximate contribution |
|---|---|---|
| `github.com/docker/docker` (moby) | Logger interface + implementations | Large: includes container runtime types, API types, networking types |
| `github.com/containerd/containerd/v2` | Transitive from moby | Large: full containerd client + server types |
| `google.golang.org/grpc` | Transitive from containerd | Large: HTTP/2 transport, protobuf codegen |
| `google.golang.org/protobuf` | Transitive from gRPC | Medium: protobuf runtime + reflection |
| `go.opentelemetry.io/otel` + related | Transitive from containerd/moby | Medium: tracing + metrics SDK |
| `github.com/prometheus/client_golang` | Transitive from moby metrics | Medium: metrics registry + HTTP handler |
| `github.com/spf13/viper` | Direct: config parsing | Medium: pulls in fsnotify, afero, toml, mapstructure, yaml |

The shim-logger only uses a small fraction of the code in these dependencies. For example, it
uses `dockerawslogs.New()` to create a CloudWatch logger, but the import chain pulls in the
entire Docker daemon logger registry, which includes code for all log drivers (json-file,
syslog, journald, gelf, etc.).

Two levels of optimization are possible:

#### Level 1: Replace viper with direct flag parsing (moderate effort)

The shim-logger uses `spf13/viper` solely to bind `pflag` flags and read their values via
`viper.GetString()`, `viper.GetBool()`, etc. It does not use viper's file-based configuration,
remote config, or watch features. Replacing viper with direct `pflag` value access
(`pflag.CommandLine.GetString()`) would remove viper and its transitive dependencies: `fsnotify`,
`go-toml`, `mapstructure`, `afero`, `cast`, `gotenv`, `sagikazarmark/locafero`, and
`sourcegraph/conc`.

#### Level 2: Replace Docker logger imports with direct SDK calls (high effort)

For the awslogs driver specifically, the shim-logger could call the CloudWatch Logs
`PutLogEvents` API directly using the already-imported `aws-sdk-go-v2/service/cloudwatchlogs`,
bypassing the Docker awslogs driver entirely. This would remove the dependency on
`github.com/docker/docker` and its massive transitive closure. However, this means reimplementing
the batching logic (grouping events into batches that fit within CloudWatch's 1 MiB batch size
limit and 10,000 event count limit), retry logic, and sequence token management that the Docker
driver currently handles.

### Trade-offs

- **Level 1 (replace viper)**: Low risk. The shim-logger's flag parsing is straightforward and
  does not benefit from viper's advanced features. The main cost is the code change itself
  (replacing `viper.GetString(key)` calls with `pflag` equivalents throughout `args.go` and
  `init.go`). There is a subtle behavior difference: viper supports environment variable binding
  (used for `SPLUNK_TOKEN` via `viper.BindEnv`), which would need to be handled manually with
  `os.Getenv` or pflag's own env binding.
- **Level 2 (replace Docker loggers)**: High risk and high effort. The Docker awslogs driver has
  years of production hardening, including edge case handling for CloudWatch API throttling,
  sequence token conflicts, log group/stream creation, multiline pattern matching, and EMF format
  detection. Reimplementing this is a significant undertaking with a large surface area for bugs.
  The splunk and fluentd drivers have similar complexity. This should only be considered if the
  memory savings from Level 1 and the other recommendations are insufficient.
- **Customer impact**: None for Level 1. For Level 2, any behavioral differences in batching,
  retry, or error handling could affect log delivery reliability. Extensive integration testing
  against real CloudWatch, Splunk, and Fluentd endpoints would be required.

---

## Recommendation 8: Add Dropped Message Counter

### Impact

Provides observability into log loss in non-blocking mode. This does not reduce memory usage but
gives customers and operators the information needed to right-size their buffer configuration.

### Why This Works

The current `Enqueue` method silently drops messages when the buffer is full:

```go
if len(b.queue) > 0 && b.curSizeInBytes+lineSizeInBytes > b.maxSizeInBytes {
    b.wait.Signal()
    return nil  // message silently dropped
}
```

There is no counter, log, or metric emitted when this happens. The ECS documentation states:
"When the buffer fills up, further logs cannot be stored. Logs that cannot be stored are lost."
But neither the customer nor the operator has any way to know *how many* logs were lost or *when*
the drops occurred.

Adding an atomic counter is trivial:

```go
var messagesDropped uint64

// In Enqueue, when dropping:
atomic.AddUint64(&messagesDropped, 1)

// In the existing startTracingLogRouting ticker (common.go), emit the count:
previousDropped := atomic.SwapUint64(&messagesDropped, 0)
debug.SendEventsToLog(containerID,
    fmt.Sprintf("Dropped %d messages in the last minute due to full buffer", previousDropped),
    debug.INFO, 0)
```

### Trade-offs

- **Minimal overhead**: A single `atomic.AddUint64` per dropped message. This is a lock-free
  operation that costs a few nanoseconds.
- **Log noise**: Emitting a log line every minute about dropped messages adds one line per minute
  to the shim-logger's own debug output. This is only visible in the shim-logger's log file, not
  in the customer's container logs.
- **No customer impact on log delivery**: This is purely additive observability. No existing
  behavior changes.

---

## Testing Strategy: Scenario-Based Memory Benchmarks

Rather than per-recommendation unit tests, we use a single suite of scenario-based benchmarks
that represent common and extreme real-world usage patterns. Each scenario measures `HeapInuse`
(via `runtime.ReadMemStats`) after a deterministic workload completes. The measured values are
recorded as thresholds in the test file.

The workflow is:

1. Write the scenarios and run them against the current (unoptimized) code.
2. Record the observed `HeapInuse` for each scenario as the initial baseline thresholds.
3. As each recommendation is implemented, re-run the suite. The measured values should decrease.
4. Update the thresholds downward to lock in the improvement and prevent regressions.

This gives a single, consistent regression gate that tracks cumulative progress across all
optimizations.

### Location

`logger/memory_benchmark_test.go` with a `//go:build unit` tag so it runs as part of the
existing `make test-unit` target.

### Scenarios

The scenarios are chosen to isolate different memory contributors and match the configurations
from the original benchmarking exercise:

| Scenario | Mode | Line Size | Buffer Size | What It Stresses |
|---|---|---|---|---|
| S1: Small lines, blocking | blocking | 100 B | n/a | Baseline: read buffers + message allocation + stream overhead |
| S2: Large lines, blocking | blocking | 62 KiB | n/a | Read buffer sizing (awslogs 256 KiB buffer) |
| S3: Small lines, non-blocking, default buffer | non-blocking | 100 B | 10 MiB | Ring buffer with many small messages (per-message overhead dominance) |
| S4: Small lines, non-blocking, large buffer | non-blocking | 100 B | 256 MiB | Ring buffer scaling — exposes pointer retention and GC pressure |
| S5: Large lines, non-blocking, default buffer | non-blocking | 62 KiB | 10 MiB | Few large messages in ring buffer |
| S6: Large lines, non-blocking, large buffer | non-blocking | 62 KiB | 512 MiB | Extreme case from original benchmarking (1.6+ GiB observed) |

### Test Structure

Each scenario follows the same pattern:

```go
func TestMemoryScenario_S3_SmallLinesNonBlockingDefault(t *testing.T) {
    const (
        lineSize      = 100
        maxBufferSize = 10 * 1024 * 1024  // 10 MiB
        numMessages   = 50_000            // enough to fill and churn the buffer
        // Baseline threshold — measured against unoptimized code.
        // Update this value downward as optimizations land.
        maxHeapInuse  = 45 * 1024 * 1024  // 45 MiB (example initial baseline)
    )

    // 1. Create a mock destination that accepts and discards messages.
    mockStream := &discardClient{}

    // 2. Create the logger in non-blocking mode.
    l, _ := logger.NewLogger(
        logger.WithStdout(/* pipe with test data */),
        logger.WithStderr(/* empty pipe */),
        logger.WithStream(mockStream),
        logger.WithInfo(logger.NewInfo("test-id", "test-name")),
    )
    bl := logger.NewBufferedLogger(l, logger.DefaultBufSizeInBytes, maxBufferSize, "test-id")

    // 3. Feed numMessages through the ring buffer.
    //    Use a bytes.Buffer as a mock pipe that emits `numMessages` lines of `lineSize` bytes.
    var pipe bytes.Buffer
    line := strings.Repeat("a", lineSize)
    for i := 0; i < numMessages; i++ {
        pipe.WriteString(line)
        pipe.WriteByte('\n')
    }

    // 4. Run the logger (Start reads from pipe, enqueues to ring buffer,
    //    consumer goroutine dequeues and calls mockStream.Log).
    ctx, cancel := context.WithCancel(context.Background())
    cleanupTime := 1 * time.Second
    go func() {
        // Wait for pipe to drain, then cancel.
        time.Sleep(5 * time.Second)
        cancel()
    }()
    _ = bl.Start(ctx, &cleanupTime, func() error { return nil })

    // 5. Force GC and measure.
    runtime.GC()
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    // 6. Assert against threshold.
    if memStats.HeapInuse > maxHeapInuse {
        t.Errorf("HeapInuse = %d bytes (%d MiB), exceeds threshold of %d bytes (%d MiB)",
            memStats.HeapInuse, memStats.HeapInuse/1024/1024,
            maxHeapInuse, maxHeapInuse/1024/1024)
    }

    // 7. Log the actual value so it can be used to update the threshold.
    t.Logf("HeapInuse = %d bytes (%d MiB)", memStats.HeapInuse, memStats.HeapInuse/1024/1024)
}
```

### How Thresholds Evolve

The thresholds are updated manually after each optimization lands. The process:

1. Implement the optimization (e.g., Recommendation 1: strip symbols).
2. Run the full scenario suite: `go test -tags unit -run TestMemoryScenario -v ./logger/`
3. Each test logs its actual `HeapInuse`. If the value decreased, update `maxHeapInuse` in the
   test to the new (lower) value plus a small margin (e.g., 10%) for noise tolerance.
4. Commit the updated thresholds alongside the optimization code.

Example progression for S3 (small lines, non-blocking, 10 MiB buffer):

| Change | Measured HeapInuse | New Threshold |
|---|---|---|
| Baseline (no changes) | ~42 MiB | 45 MiB |
| After Rec 1 (strip symbols) | ~42 MiB (no heap effect) | 45 MiB (unchanged) |
| After Rec 6 (GOGC=50) | ~32 MiB | 35 MiB |
| After Rec 3 (GOMEMLIMIT) | ~25 MiB | 28 MiB |
| After Rec 2 (nil-out pointers) | ~18 MiB | 20 MiB |
| After Rec 4 (overhead accounting) | ~12 MiB | 14 MiB |

Note: The values above are illustrative. The actual baseline must be measured. The key property
is that thresholds only move downward — any change that increases `HeapInuse` beyond the current
threshold will fail the test, catching regressions.

### Binary Size Scenario

In addition to the heap scenarios, one scenario tracks binary size. This is not a
`runtime.ReadMemStats` test but a build-and-measure step:

```go
func TestBinarySize(t *testing.T) {
    const (
        // Baseline threshold — measured against current build.
        // Update downward as ldflags and dependency reductions land.
        maxBinarySize = 40 * 1024 * 1024  // 40 MiB (example, calibrate after first build)
    )

    // Build the binary using the same flags as the Makefile.
    cmd := exec.Command("go", "build", "-o", "/tmp/shim-logger-test-binary", ".")
    cmd.Dir = ".."  // workspace root
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("build failed: %s\n%s", err, out)
    }
    defer os.Remove("/tmp/shim-logger-test-binary")

    info, err := os.Stat("/tmp/shim-logger-test-binary")
    if err != nil {
        t.Fatal(err)
    }

    if info.Size() > maxBinarySize {
        t.Errorf("binary size = %d bytes (%d MiB), exceeds threshold of %d bytes (%d MiB)",
            info.Size(), info.Size()/1024/1024,
            maxBinarySize, int64(maxBinarySize)/1024/1024)
    }
    t.Logf("binary size = %d bytes (%d MiB)", info.Size(), info.Size()/1024/1024)
}
```

Example progression:

| Change | Measured Size | New Threshold |
|---|---|---|
| Baseline (no ldflags) | ~35 MiB | 40 MiB |
| After Rec 1 (strip symbols) | ~25 MiB | 28 MiB |
| After Rec 7 Level 1 (remove viper) | ~23 MiB | 25 MiB |
| After Rec 7 Level 2 (remove Docker loggers) | ~15 MiB | 17 MiB |

### Why This Approach

- **One suite, all optimizations**: Instead of per-recommendation tests that each test one
  mechanism in isolation, the scenarios test the observable outcome (heap size) that customers
  care about. Any optimization that reduces heap usage — whether it's GC tuning, pointer
  nil-out, or overhead accounting — shows up as a lower number in the same test.
- **Ratcheting prevents regressions**: Because thresholds only move down, a future change that
  adds a dependency or introduces a memory leak will fail the scenario that it affects.
- **Reproducible**: The scenarios use deterministic workloads (fixed number of messages, fixed
  sizes) with `runtime.GC()` before measurement, minimizing noise. The 10% margin on thresholds
  accounts for allocator non-determinism and GC timing.
- **Low maintenance**: Adding a new scenario is just adding a new test function with the same
  pattern. Updating a threshold is a one-line change.

---

## Summary and Prioritization

| # | Recommendation | Effort | RSS Impact | Risk | Status |
|---|---|---|---|---|---|
| 1 | Strip debug symbols (`-s -w`) | Trivial | -33% binary size (26→17.5 MiB) | None | ✅ Implemented |
| 2 | Nil out dequeued pointers | Low | 107.6→0.7 MiB retained after drain | None | ✅ Implemented |
| 3 | Set GOMEMLIMIT | Low | High (caps worst-case RSS) | GC thrashing under sustained load | Future |
| 4 | Per-message overhead accounting | Low | Moderate (accurate tracking) | More log drops for small messages | Future |
| 5 | Dynamic read buffer for awslogs | Medium | Low (~480 KiB per instance) | More partial messages for large lines | Future |
| 6 | GOGC tuning | Trivial | Moderate (peak RSS) | CPU/memory tradeoff | Future |
| 7 | Dependency reduction | High | High (baseline RSS) | Reimplementation risk (Level 2) | Future |
| 8 | Dropped message counter | Trivial | None (observability) | None | Future |

Recommendations 1 and 2 are shipped. The remaining recommendations should be revisited after
measuring the production impact of these changes. If peak RSS during sustained throughput
(the S3/S5 scenarios) remains problematic, Recommendations 3 and 6 are the next candidates —
but they trade CPU for memory and should be validated against real workload profiles before
deployment.
