//go:build windows
// +build windows

// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cihub/seelog"
)

const (
	// Maximum file size in MB
	MAX_FILE_SIZE float64 = 10
	// Maximum number of files generated in log rotation
	MAX_ROLLS int = 24
	// Messages will be in the format "time=... level=... msg=..."
	SEE_LOG_CONFIG_TEMPLATE string = `
<seelog type="asyncloop">
	<outputs formatid="main">
		<rollingfile filename="%s" type="size"
			maxsize="%d" archivetype="none" maxrolls="%d" />
	</outputs>
	<formats>
		<format id="main" format="time=%%Date/%%Time level=%%LEV msg=%%Msg%%n"/>
	</formats>
</seelog>`
)

var (
	envProgramData string = defaultIfBlank(os.Getenv("ProgramData"), `C:\ProgramData`)
	// logFileDir is the default path for logs when log-file-dir option is not set
	logFileDir    string = filepath.Join(envProgramData, "Amazon/ECS/log/shim-logger")
	containerName string
	fileLogger    seelog.LoggerInterface
	fileLoggerMu  sync.Mutex
)

func FlushLog() {
	fileLoggerMu.Lock()
	defer fileLoggerMu.Unlock()
	fileLogger.Flush()
}

func SendEventsToLog(logfileNameId string, msg string, msgType string, delay time.Duration) {
	filename := fmt.Sprintf("%s-%s.log", containerName, logfileNameId)
	file := filepath.Join(logFileDir, filename)
	configStr := fmt.Sprintf(SEE_LOG_CONFIG_TEMPLATE, file, int(MAX_FILE_SIZE*1000000), MAX_ROLLS)

	fileLoggerMu.Lock()
	fileLogger, _ = seelog.LoggerFromConfigAsString(configStr)
	logger := fileLogger
	fileLoggerMu.Unlock()

	switch msgType {
	case "err":
		logger.Error(msg)
	case "info":
		logger.Info(msg)
	case "debug":
		logger.Debug(msg)
	}
	time.Sleep(delay * time.Second)
}

func SetLogFilePath(logFlag, contName string) error {
	containerName = contName
	logFileDir = logFlag
	return nil
}

// Not implemented in Windows
func StartStackTraceHandler() {}

func defaultIfBlank(str, defaultValue string) string {
	if len(str) == 0 {
		return defaultValue
	}
	return str
}
