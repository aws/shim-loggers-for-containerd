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

func WithStream(stream client) LoggerOpt {
	return func(l *Logger) {
		l.Stream = stream
	}
}
