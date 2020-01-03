package logger

import "time"

const (
	DaemonName = "shim-loggers-for-containerd"

	// Define the retry parameters for retrying driving logs to destination
	LogRetryMaxAttempts = 3
	LogRetryMinBackoff  = 500 * time.Millisecond
	LogRetryMaxBackoff  = 1 * time.Second
	LogRetryJitter      = 0.3
	LogRetryMultiple    = 2
)

type GlobalArgs struct {
	// Required arguments
	ContainerID   string
	ContainerName string
	LogDriver     string

	// Optional arguments
	Mode          string
	MaxBufferSize int
}

// Interface for all log drivers
type LogDriver interface {
	// Start functions starts driving container debug to destination
	Start(func() error) error
}
