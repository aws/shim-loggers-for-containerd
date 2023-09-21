// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"github.com/containerd/containerd/cio"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	splunkTokenKey              = "--splunk-token"
	splunkUrlkey                = "--splunk-url"
	splunkInsecureskipverifyKey = "--splunk-insecureskipverify"
	testSplunkUrl               = "https://localhost:8089"
)

var testSplunk = func(token string) {
	// These tests are run in serial because we only define one log driver instance.
	ginkgo.Describe("splunk shim logger", ginkgo.Serial, func() {
		gomega.Expect(token).ShouldNot(gomega.BeEmpty())
		ginkgo.It("should send logs to splunk log driver", func() {
			args := map[string]string{
				logDriverTypeKey:            splunkDriverName,
				containerIdKey:              testContainerId,
				containerNameKey:            testContainerName,
				splunkTokenKey:              token,
				splunkUrlkey:                testSplunkUrl,
				splunkInsecureskipverifyKey: "true",
			}
			creator := cio.BinaryIO(*Binary, args)
			sendTestLogByContainerd(creator, testLog)
			// TODO: Validate logs in Splunk local. https://github.com/aws/shim-loggers-for-containerd/issues/74
		})
	})
}
