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
	"io"

	dockerlogger "github.com/docker/docker/daemon/logger"
)

// LoggerOpt is a type of function that is used to update the values
// of fields in LoggerArgs. Fields supported to be modified are
// logger info, stdout and stderr.
type LoggerOpt func(*Logger)

// InfoOpt is a type of function that is used to update the values
// of fields in logger info for each driver. Field supported to be
// modified is config.
type InfoOpt func(*dockerlogger.Info)

// WithConfig sets logger config of logger info
func WithConfig(m map[string]string) InfoOpt {
	return func(info *dockerlogger.Info) {
		info.Config = m
	}
}

// WithStdout sets log driver's stdout pipe
func WithStdout(stdout io.Reader) LoggerOpt {
	return func(l *Logger) {
		l.Stdout = stdout
	}
}

// WithStderr sets log driver's stderr pipe
func WithStderr(stderr io.Reader) LoggerOpt {
	return func(l *Logger) {
		l.Stderr = stderr
	}
}

// WithInfo sets log driver's info
func WithInfo(info *dockerlogger.Info) LoggerOpt {
	return func(l *Logger) {
		l.Info = info
	}
}

func WithStream(stream Client) LoggerOpt {
	return func(l *Logger) {
		l.Stream = stream
	}
}
