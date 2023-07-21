// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"io"

	dockerlogger "github.com/docker/docker/daemon/logger"
)

// Opt is a type of function that is used to update the values
// of fields in LoggerArgs. Fields supported to be modified are
// logger info, stdout and stderr.
type Opt func(*Logger)

// InfoOpt is a type of function that is used to update the values
// of fields in logger info for each driver. Field supported to be
// modified is config.
type InfoOpt func(*dockerlogger.Info)

// WithConfig sets logger config of logger info.
func WithConfig(m map[string]string) InfoOpt {
	return func(info *dockerlogger.Info) {
		info.Config = m
	}
}

// WithStdout sets log driver's stdout pipe.
func WithStdout(stdout io.Reader) Opt {
	return func(l *Logger) {
		l.Stdout = stdout
	}
}

// WithStderr sets log driver's stderr pipe.
func WithStderr(stderr io.Reader) Opt {
	return func(l *Logger) {
		l.Stderr = stderr
	}
}

// WithInfo sets log driver's info.
func WithInfo(info *dockerlogger.Info) Opt {
	return func(l *Logger) {
		l.Info = info
	}
}

// WithStream sets the actual stream of log driver.
func WithStream(stream Client) Opt {
	return func(l *Logger) {
		l.Stream = stream
	}
}

// WithBufferSizeInBytes sets the buffer size of log driver.
func WithBufferSizeInBytes(size int) Opt {
	return func(l *Logger) {
		l.bufferSizeInBytes = size
	}
}

// WithMaxReadBytes sets how many bytes will be read from container
// pipe per iteration.
func WithMaxReadBytes(size int) Opt {
	return func(l *Logger) {
		l.maxReadBytes = size
	}
}
