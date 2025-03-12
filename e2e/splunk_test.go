// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build e2e

package e2e

import (
	"github.com/containerd/containerd/cio"
	"github.com/google/uuid"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	splunkTokenKey              = "--splunk-token" //nolint:gosec // no real credential
	splunkURLkey                = "--splunk-url"
	splunkInsecureskipverifyKey = "--splunk-insecureskipverify"
	testSplunkURL               = "https://localhost:8089"
)

var testSplunk = func(token string) {
	// These tests are run in serial because we only define one log driver instance.
	ginkgo.Describe("splunk shim logger", ginkgo.Serial, func() {
		gomega.Expect(token).ShouldNot(gomega.BeEmpty())
		ginkgo.It("should send logs to splunk log driver", func() {
			testLog := testLogPrefix + uuid.New().String()
			args := map[string]string{
				LogDriverTypeKey:            SplunkDriverName,
				ContainerIDKey:              TestContainerID,
				ContainerNameKey:            TestContainerName,
				splunkTokenKey:              token,
				splunkURLkey:                testSplunkURL,
				splunkInsecureskipverifyKey: "true",
			}
			creator := cio.BinaryIO(*Binary, args)
			err := SendTestLogByContainerd(creator, testLog)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			// TODO: Validate logs in Splunk local. https://github.com/aws/shim-loggers-for-containerd/issues/74
		})
	})
}
