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

package awslogs

import (
	"github.com/aws/shim-loggers-for-containerd/logger"

	dockerlogger "github.com/docker/docker/daemon/logger"
)

// WithRegion sets awslogs region of logger info
func WithRegion(region string) logger.InfoOpt {
	return func(info *dockerlogger.Info) {
		info.Config[RegionKey] = region
	}
}
