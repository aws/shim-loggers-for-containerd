// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build e2e

package e2e

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/cio"
	"github.com/google/uuid"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	jsonFileLogPathKey  = "--log-path"
	jsonFileMaxSizeKey  = "--max-size"
	jsonFileMaxFileKey  = "--max-file"
	jsonFileCompressKey = "--compress"
	jsonFileTagKey      = "--json-file-tag"

	// jsonFileLogDir is the per-test output directory where the shim-logger writes the
	// json-file output for the test container. Files are named "<container-id>-json.log"
	// plus rotated siblings.
	//
	// Relative path used here is resolved to an absolute path at test-setup time —
	// the shim-logger is invoked by containerd and inherits a different CWD than the
	// test runner, so a relative path would fail to resolve in the shim's context.
	jsonFileLogDir   = "../jsonfile-logs"
	jsonFileLogName  = TestContainerID + "-json.log"
	expectedFileMode = 0o640
)

// jsonFileEnvelope models the per-line shape moby's jsonfilelog produces:
//
//	{"log":"...","stream":"stdout|stderr","time":"<RFC3339Nano>"}
type jsonFileEnvelope struct {
	Log    string `json:"log"`
	Stream string `json:"stream"`
	Time   string `json:"time"`
}

var testJSONFile = func() {
	// absLogDir is the absolute path to jsonFileLogDir. Computed once at suite-build time
	// so all sub-specs see the same value. The shim-logger runs with a different CWD than
	// the test runner (it's launched by containerd), so the path passed via --log-path
	// must be absolute.
	absLogDir, err := filepath.Abs(jsonFileLogDir)
	if err != nil {
		panic(fmt.Sprintf("failed to resolve absolute path for %q: %v", jsonFileLogDir, err))
	}
	absLogFile := filepath.Join(absLogDir, jsonFileLogName)

	// Tests are serial because we share a single output dir and a single TestContainerID.
	ginkgo.Describe("json-file shim logger", ginkgo.Serial, func() {
		ginkgo.BeforeEach(func() {
			// Ensure a fresh output directory for each test.
			gomega.Expect(os.RemoveAll(absLogDir)).Should(gomega.Succeed())
			gomega.Expect(os.MkdirAll(absLogDir, 0o750)).Should(gomega.Succeed())
		})

		ginkgo.It("writes Docker-envelope JSON for each output line", func() {
			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:   JSONFileDriverName,
				ContainerIDKey:     TestContainerID,
				ContainerNameKey:   TestContainerName,
				jsonFileLogPathKey: absLogFile,
			}
			creator := cio.BinaryIO(*Binary, args)

			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			lines := readEnvelopeLines(absLogFile)
			gomega.Expect(lines).ShouldNot(gomega.BeEmpty())
			// We assert the test log shows up at least once across the produced envelopes.
			// Not asserting count since printf may emit a single line or be split.
			joined := joinEnvelopeLogs(lines)
			gomega.Expect(joined).Should(gomega.ContainSubstring(testLog))
			// Each envelope's stream should be stdout (printf writes to stdout).
			for _, env := range lines {
				gomega.Expect(env.Stream).Should(gomega.Equal("stdout"))
				gomega.Expect(env.Time).ShouldNot(gomega.BeEmpty())
			}
		})

		ginkgo.It("rotates files when --max-size is exceeded", func() {
			args := map[string]string{
				LogDriverTypeKey:   JSONFileDriverName,
				ContainerIDKey:     TestContainerID,
				ContainerNameKey:   TestContainerName,
				jsonFileLogPathKey: absLogFile,
				jsonFileMaxSizeKey: "1k",
				jsonFileMaxFileKey: "3",
			}
			creator := cio.BinaryIO(*Binary, args)

			// Emit ~3 KB of distinct lines so we force at least 2 rotations with --max-size=1k.
			// Each line is ~70 bytes; 50 lines * ~70B = ~3.5 KB after envelope overhead.
			err := SendCommandByContainerd(creator, `for i in $(seq 1 50); do echo "rotation-test-line-$i-padding-padding-padding"; done`)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			files := listLogFiles(absLogDir, jsonFileLogName)
			// With --max-file=3 we expect the active file plus up to 2 rotated siblings.
			gomega.Expect(len(files)).Should(gomega.BeNumerically(">=", 2),
				"expected at least one rotated file plus the active one, got %v", files)
			gomega.Expect(len(files)).Should(gomega.BeNumerically("<=", 3),
				"max-file=3 should cap the on-disk file count at 3, got %v", files)
		})

		ginkgo.It("caps the file count at --max-file beyond the cap", func() {
			args := map[string]string{
				LogDriverTypeKey:   JSONFileDriverName,
				ContainerIDKey:     TestContainerID,
				ContainerNameKey:   TestContainerName,
				jsonFileLogPathKey: absLogFile,
				jsonFileMaxSizeKey: "1k",
				jsonFileMaxFileKey: "2",
			}
			creator := cio.BinaryIO(*Binary, args)

			// Emit a lot more output so we'd rotate well past max-file.
			err := SendCommandByContainerd(creator, `for i in $(seq 1 200); do echo "cap-test-line-$i-padding-padding-padding-padding"; done`)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			files := listLogFiles(absLogDir, jsonFileLogName)
			gomega.Expect(len(files)).Should(gomega.BeNumerically("<=", 2),
				"max-file=2 should cap the on-disk file count, got %v", files)
		})

		ginkgo.It("does not rotate when total output is below --max-size", func() {
			args := map[string]string{
				LogDriverTypeKey:   JSONFileDriverName,
				ContainerIDKey:     TestContainerID,
				ContainerNameKey:   TestContainerName,
				jsonFileLogPathKey: absLogFile,
				jsonFileMaxSizeKey: "1m",
				jsonFileMaxFileKey: "5",
			}
			creator := cio.BinaryIO(*Binary, args)

			testLog := testLogPrefix + uuid.New().String()
			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			files := listLogFiles(absLogDir, jsonFileLogName)
			gomega.Expect(files).Should(gomega.HaveLen(1),
				"a single small printf should not trigger rotation, got %v", files)
		})

		ginkgo.It("creates the output file with mode 0640", func() {
			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:   JSONFileDriverName,
				ContainerIDKey:     TestContainerID,
				ContainerNameKey:   TestContainerName,
				jsonFileLogPathKey: absLogFile,
			}
			creator := cio.BinaryIO(*Binary, args)

			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			info, err := os.Stat(absLogFile)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			// File mode bits only — strip the type bits for a clean comparison.
			gomega.Expect(info.Mode().Perm()).Should(gomega.Equal(os.FileMode(expectedFileMode)),
				"expected mode %#o, got %#o", expectedFileMode, info.Mode().Perm())
		})

		ginkgo.It("survives buffer pressure under non-blocking mode", func() {
			// Non-blocking mode wraps the json-file writer in shim-logger's
			// ringBuffer. With a tiny --max-buffer-size and a fast producer,
			// the buffer fills and lines are dropped — that's moby/shim-logger
			// contract. This spec doesn't assert *which* lines drop (timing-
			// dependent); it asserts that:
			//   (1) the shim-logger doesn't crash,
			//   (2) the container task exits cleanly,
			//   (3) every line that *did* land in the file is a valid Docker
			//       envelope (no torn writes, no partial JSON).
			//
			// This exercises the NonBlockingMode wrapping in
			// jsonfile.RunLogDriver, which is new in this PR.
			args := map[string]string{
				LogDriverTypeKey:    JSONFileDriverName,
				ContainerIDKey:      TestContainerID,
				ContainerNameKey:    TestContainerName,
				jsonFileLogPathKey:  absLogFile,
				"--mode":            "non-blocking",
				"--max-buffer-size": "64k",
			}
			creator := cio.BinaryIO(*Binary, args)

			// Emit ~200 KiB of output as fast as the shell can produce it.
			// Each line is ~50 bytes; 4000 lines * ~50B = ~200 KiB. A 64 KiB
			// buffer plus the writer draining means many lines will queue
			// and some will be dropped.
			err := SendCommandByContainerd(creator,
				`for i in $(seq 1 4000); do echo "non-blocking-line-$i"; done`)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred(),
				"shim-logger should exit cleanly even when the buffer overflows")

			// Whatever lines made it through must be well-formed envelopes.
			// readEnvelopeLines fails fast on any unparseable line.
			lines := readEnvelopeLines(absLogFile)
			gomega.Expect(lines).ShouldNot(gomega.BeEmpty(),
				"at least some lines should make it through even under buffer pressure")
			for _, env := range lines {
				gomega.Expect(env.Stream).Should(gomega.Equal("stdout"))
				gomega.Expect(env.Time).ShouldNot(gomega.BeEmpty())
			}
		})

		ginkgo.It("the binary fails fast when --log-path is missing", func() {
			args := map[string]string{
				LogDriverTypeKey: JSONFileDriverName,
				ContainerIDKey:   TestContainerID,
				ContainerNameKey: TestContainerName,
				// Intentionally no --log-path.
			}
			creator := cio.BinaryIO(*Binary, args)

			err := SendTestLogByContainerd(creator, testLogPrefix+uuid.New().String())
			gomega.Expect(err).Should(gomega.HaveOccurred(),
				"missing --log-path should cause the shim-logger to fail and the container task to error out")
		})
	})
}

// readEnvelopeLines parses every JSON-line in the given file as a jsonFileEnvelope.
// Empty lines are skipped. The test fails fast on any parse error.
func readEnvelopeLines(path string) []jsonFileEnvelope {
	file, err := os.Open(path) //nolint:gosec // testing only
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer file.Close() //nolint:errcheck // closing the file

	var envelopes []jsonFileEnvelope
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var env jsonFileEnvelope
		err := json.Unmarshal([]byte(line), &env)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred(), "line %q failed to parse as JSON envelope", line)
		envelopes = append(envelopes, env)
	}
	gomega.Expect(scanner.Err()).ShouldNot(gomega.HaveOccurred())
	return envelopes
}

// joinEnvelopeLogs concatenates the .Log fields of every envelope in order.
// Useful for asserting that a multi-line printf output appears across envelopes.
func joinEnvelopeLogs(envelopes []jsonFileEnvelope) string {
	var b strings.Builder
	for _, e := range envelopes {
		b.WriteString(e.Log)
	}
	return b.String()
}

// listLogFiles returns the active log file plus any rotated siblings (e.g.,
// "<name>", "<name>.1", "<name>.2", ...).
func listLogFiles(dir, baseName string) []string {
	entries, err := os.ReadDir(dir)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	var matches []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == baseName || strings.HasPrefix(name, baseName+".") {
			matches = append(matches, name)
		}
	}
	if len(matches) == 0 {
		// Provide a helpful failure message.
		var present []string
		for _, e := range entries {
			present = append(present, e.Name())
		}
		gomega.Expect(matches).ShouldNot(gomega.BeEmpty(),
			"no files matching base %q found in %q (present: %v)", baseName, dir, present)
	}
	return matches
}
