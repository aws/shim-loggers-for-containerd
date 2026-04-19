// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build !linux && !windows

package logger

import "errors"

// UID not supported in Windows.
func setUID(_ int) error {
	return errors.New("UID not supported")
}

// GID not supported in Windows.
func setGID(_ int) error {
	return errors.New("GID not supported")
}
