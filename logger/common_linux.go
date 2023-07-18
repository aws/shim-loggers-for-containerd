//go:build !windows
// +build !windows

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

package logger

import (
	"fmt"
	"syscall"

	"github.com/aws/shim-loggers-for-containerd/debug"
)

// setUID sets UID of current goroutine/process.
// If you are building with go version includes the following commit, this syscall would apply
// to current process, otherwise it would only apply to current goroutine.
// Commit: https://github.com/golang/go/commit/d1b1145cace8b968307f9311ff611e4bb810710c
func setUID(id int) error {
	err := syscall.Setuid(id)
	if err != nil {
		return fmt.Errorf("unable to set uid: %w", err)
	}

	// Check if uid set correctly
	u := syscall.Getuid()
	if u != id {
		return fmt.Errorf("want uid %d, but get uid %d", id, u)
	}
	debug.SendEventsToLog(DaemonName,
		fmt.Sprintf("Set uid: %d", u),
		debug.INFO, 1)

	return nil
}

// setGID sets GID of current goroutine/process.
// If you are building with go version includes the following commit, this syscall would apply
// to current process, otherwise it would only apply to current goroutine.
// Commit: https://github.com/golang/go/commit/d1b1145cace8b968307f9311ff611e4bb810710c
func setGID(id int) error {
	err := syscall.Setgid(id)
	if err != nil {
		return fmt.Errorf("unable to set gid: %w", err)
	}

	// Check if gid set correctly
	g := syscall.Getgid()
	if g != id {
		return fmt.Errorf("want gid %d, but get gid %d", id, g)
	}
	debug.SendEventsToLog(DaemonName,
		fmt.Sprintf("Set gid %d", g),
		debug.INFO, 1)

	return nil
}
