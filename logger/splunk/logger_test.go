// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package splunk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testToken              = "testToken"
	testURL                = "localhost:8000"
	testSource             = "httpSource"
	testSourceType         = "http"
	testIndex              = "main"
	testCapath             = "/root"
	testCaname             = "localhost"
	testInsecureskipverify = "true"
	testFormat             = "inline"
	testVerifyConnection   = "true"
	testGzip               = "true"
	testGzipLevel          = "0"
	testSplunkTag          = "tag"
	testLabels             = "label0, label1"
	testEnv                = "envVar0, envVar1"
	testEnvRegex           = "envVar*"
)

var (
	testArg = &Args{
		Token:              testToken,
		URL:                testURL,
		Source:             testSource,
		Sourcetype:         testSourceType,
		Index:              testIndex,
		Capath:             testCapath,
		Caname:             testCaname,
		Insecureskipverify: testInsecureskipverify,
		Format:             testFormat,
		VerifyConnection:   testVerifyConnection,
		Gzip:               testGzip,
		GzipLevel:          testGzipLevel,
		Tag:                testSplunkTag,
		Labels:             testLabels,
		Env:                testEnv,
		EnvRegex:           testEnvRegex,
	}
)

// TestGetSplunkConfig tests if all arguments are converted in correctly
// for splunk driver configuration
func TestGetSplunkConfig(t *testing.T) {
	expectedConfig := map[string]string{
		TokenKey:              testToken,
		URLKey:                testURL,
		SourceKey:             testSource,
		SourcetypeKey:         testSourceType,
		IndexKey:              testIndex,
		CapathKey:             testCapath,
		CanameKey:             testCaname,
		InsecureskipverifyKey: testInsecureskipverify,
		FormatKey:             testFormat,
		VerifyConnectionKey:   testVerifyConnection,
		GzipKey:               testGzip,
		GzipLevelKey:          testGzipLevel,
		tagKey:                testSplunkTag,
		LabelsKey:             testLabels,
		EnvKey:                testEnv,
		EnvRegexKey:           testEnvRegex,
	}
	config := getSplunkConfig(testArg)
	require.Equal(t, expectedConfig, config)
}
