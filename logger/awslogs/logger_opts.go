// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

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
