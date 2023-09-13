// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import "github.com/onsi/ginkgo/v2"

var testFluentd = func() {
	// These tests are run in serial because we only define one log driver instance.
	ginkgo.Describe("fluentd shim logger", ginkgo.Serial, func() {
		// TODO: add test cases
	})
}
