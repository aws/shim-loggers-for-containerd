//go:build windows
// +build windows

// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"errors"
)

// UID not supported in Windows
func setUID(id int) error {
	return errors.New("UID not supported in Windows")
}

// GID not supported in Windows
func setGID(id int) error {
	return errors.New("GID not supported in Windows")
}
