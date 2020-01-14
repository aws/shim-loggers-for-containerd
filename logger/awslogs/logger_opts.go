package awslogs

import (
	"io"

	"github.com/docker/docker/daemon/logger"
)

// LoggerOpt is a type of function that is used to update the values
// of fields in LoggerArgs. Fields supported to be modified are
// logger info, stdout and stderr.
type LoggerOpt func(*logDriver)

// InfoOpt is a type of function that is used to update the values
// of fields in logger info. Fields supported to be modified are
// awslogs config, and region.
type InfoOpt func(*logger.Info)

// WithConfig sets logger config of logger info
func WithConfig(m map[string]string) InfoOpt {
	return func(info *logger.Info) {
		info.Config = m
	}
}

// WithRegion sets awslogs region of logger info
func WithRegion(region string) InfoOpt {
	return func(info *logger.Info) {
		info.Config[RegionKey] = region
	}
}

// WithStdout sets log driver's stdout pipe
func WithStdout(stdout io.Reader) LoggerOpt {
	return func(l *logDriver) {
		l.stdout = stdout
	}
}

// WithStderr sets log driver's stderr pipe
func WithStderr(stderr io.Reader) LoggerOpt {
	return func(l *logDriver) {
		l.stderr = stderr
	}
}

// WithInfo sets log driver's info
func WithInfo(info *logger.Info) LoggerOpt {
	return func(l *logDriver) {
		l.info = info
	}
}
