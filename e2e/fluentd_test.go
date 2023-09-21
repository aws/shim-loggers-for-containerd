// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/cio"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	fluentdLogDirName = "./../fluentd-logs"
)

var testFluentd = func() {
	// These tests are run in serial because we only define one log driver instance.
	ginkgo.Describe("fluentd shim logger", ginkgo.Serial, func() {
		ginkgo.It("should send logs to fluentd log driver", func() {
			args := map[string]string{
				logDriverTypeKey: fluentdDriverName,
				containerIdKey:   testContainerId,
				containerNameKey: testContainerName,
			}
			creator := cio.BinaryIO(*Binary, args)
			sendTestLogByContainerd(creator, testLog)
			validateTestLogsInFluentd(fluentdLogDirName, testLog)
		})
	})
}

func validateTestLogsInFluentd(dirName string, testLog string) {
	// For single test, there are 3 files in Fluentd log dir: "data.<hash>.log", "data.<hash>.log.meta" and "data.log".
	// For example: "data.b60581c99383f387cfaba1fc90272852e.log", "data.b60581c99383f387cfaba1fc90272852e.log.meta" and "data.log"
	// "data.<hash>.log" has the logs that the tests sent.
	// "data.<hash>.log" can have multiple lines of records following time sequence. Here is a sample content with 3 lines.
	// 2023-09-16T22:54:11+00:00       123456789012    {"source":"stdout","log":"test-e2e-log","container_id":"123456789012","container_name":"test-container"}
	// 2023-09-16T22:54:30+00:00       123456789012    {"container_id":"123456789012","container_name":"test-container","source":"stdout","log":"test-e2e-log"}
	// 2023-09-16T22:56:17+00:00       123456789012    {"container_id":"123456789012","container_name":"test-container","source":"stdout","log":"test-e2e-log"}
	// The following steps retrieves the "log" field of the third string parsed by tab of the last line to validate the tests sent.
	var fileName string
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), "data.") && strings.HasSuffix(info.Name(), ".log") && info.Name() != "data.log" {
			fileName = path
		}
		return nil
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(fileName).ShouldNot(gomega.Equal(""))
	file, err := os.Open(fileName)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer file.Close()
	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	contentParts := strings.Split(lastLine, "\t")
	gomega.Expect(len(contentParts)).Should(gomega.Equal(3))
	var logContent map[string]string
	err = json.Unmarshal([]byte(contentParts[2]), &logContent)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(logContent["log"]).Should(gomega.Equal(testLog))
}
