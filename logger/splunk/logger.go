package splunk

import (
	"context"

	"github.com/containerd/containerd/runtime/v2/logging"
	"github.com/coreos/go-systemd/journal"
	dockersplunk "github.com/docker/docker/daemon/logger/splunk"
	"github.com/pkg/errors"

	"github.com/aws/shim-loggers-for-containerd/debug"
	"github.com/aws/shim-loggers-for-containerd/logger"
)

// splunk driver argument keys
const (
	// Required
	TokenKey = "splunk-token"
	URLKey   = "splunk-url"

	// Optional
	SourceKey             = "splunk-source"
	SourcetypeKey         = "splunk-sourcetype"
	IndexKey              = "splunk-index"
	CapathKey             = "splunk-capath"
	CanameKey             = "splunk-caname"
	InsecureskipverifyKey = "splunk-insecureskipverify"
	FormatKey             = "splunk-format"
	VerifyConnectionKey   = "splunk-verify-connection"
	GzipKey               = "splunk-gzip"
	GzipLevelKey          = "splunk-gzip-level"
	SplunkTagKey          = "splunk-tag"
	LabelsKey             = "labels"
	EnvKey                = "env"
	EnvRegexKey           = "env-regex"

	// Convert input parameter "splunk-tag" to the splunk parameter "tag"
	// This is to distinguish between the "tag" parameter from the fluentd input
	tagKey = "tag"
)

// Args represents splunk log driver arguments
type Args struct {
	Token              string
	URL                string
	Source             string
	Sourcetype         string
	Index              string
	Capath             string
	Caname             string
	Insecureskipverify string
	Format             string
	VerifyConnection   string
	Gzip               string
	GzipLevel          string
	Tag                string
	Labels             string
	Env                string
	EnvRegex           string
}

// LoggerArgs stores global logger args and splunk specific args
type LoggerArgs struct {
	globalArgs *logger.GlobalArgs
	args       *Args
}

// InitLogger initialize the input arguments
func InitLogger(globalArgs *logger.GlobalArgs, splunkArgs *Args) *LoggerArgs {
	return &LoggerArgs{
		globalArgs: globalArgs,
		args:       splunkArgs,
	}
}

// RunLogDriver initiates the splunk driver
func (la *LoggerArgs) RunLogDriver(ctx context.Context, config *logging.Config, ready func() error) error {
	loggerConfig := getSplunkConfig(la.args)
	info := logger.NewInfo(
		la.globalArgs.ContainerID,
		la.globalArgs.ContainerName,
		logger.WithConfig(loggerConfig),
	)
	stream, err := dockersplunk.New(*info)
	if err != nil {
		return errors.Wrap(err, "unable to create stream")
	}

	l, err := logger.NewLogger(
		logger.WithStdout(config.Stdout),
		logger.WithStderr(config.Stderr),
		logger.WithInfo(info),
		logger.WithStream(stream),
	)
	if err != nil {
		return errors.Wrap(err, "unable to create splunk log driver")
	}

	if la.globalArgs.Mode == logger.NonBlockingMode {
		debug.SendEventsToJournal(logger.DaemonName, "Starting non-blocking mode driver", journal.PriInfo)
		l = logger.NewBufferedLogger(l, la.globalArgs.MaxBufferSize)
	}

	// Start splunk log driver
	debug.SendEventsToJournal(logger.DaemonName, "Starting splunk driver", journal.PriInfo)
	err = l.Start(la.globalArgs.UID, la.globalArgs.GID, ready)
	if err != nil {
		return errors.Wrap(err, "failed to start splunk driver")
	}
	debug.SendEventsToJournal(logger.DaemonName, "Logging finished", journal.PriInfo)

	return nil
}

// getSplunkConfig sets values for splunk config
func getSplunkConfig(arg *Args) map[string]string {
	config := make(map[string]string)
	// Required arguments
	config[TokenKey] = arg.Token
	config[URLKey] = arg.URL
	// Optional arguments
	if arg.Source != "" {
		config[SourceKey] = arg.Source
	}
	if arg.Sourcetype != "" {
		config[SourcetypeKey] = arg.Sourcetype
	}
	if arg.Index != "" {
		config[IndexKey] = arg.Index
	}
	if arg.Capath != "" {
		config[CapathKey] = arg.Capath
	}
	if arg.Caname != "" {
		config[CanameKey] = arg.Caname
	}
	if arg.Insecureskipverify != "" {
		config[InsecureskipverifyKey] = arg.Insecureskipverify
	}
	if arg.Format != "" {
		config[FormatKey] = arg.Format
	}
	if arg.VerifyConnection != "" {
		config[VerifyConnectionKey] = arg.VerifyConnection
	}
	if arg.Gzip != "" {
		config[GzipKey] = arg.Gzip
	}
	if arg.GzipLevel != "" {
		config[GzipLevelKey] = arg.GzipLevel
	}
	if arg.Tag != "" {
		config[tagKey] = arg.Tag
	}
	if arg.Labels != "" {
		config[LabelsKey] = arg.Labels
	}
	if arg.Env != "" {
		config[EnvKey] = arg.Env
	}
	if arg.EnvRegex != "" {
		config[EnvRegexKey] = arg.EnvRegex
	}
	return config
}
