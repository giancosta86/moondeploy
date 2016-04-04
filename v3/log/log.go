/*§
  ===========================================================================
  MoonDeploy
  ===========================================================================
  Copyright (C) 2015-2016 Gianluca Costa
  ===========================================================================
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  ===========================================================================
*/

package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/giancosta86/moondeploy/v3/moonclient"
	"github.com/op/go-logging"
)

type LoggingCallback func(level logging.Level, message string)

var logger logging.Logger
var leveledBackend logging.LeveledBackend

var loggingCallback LoggingCallback = func(level logging.Level, message string) {}

func SetCallback(callback LoggingCallback) {
	loggingCallback = callback
}

func Debug(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.DEBUG) {
		message := fmt.Sprintf(format, args...)

		logger.Debug(message)

		loggingCallback(logging.DEBUG, message)
	}
}

func Info(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.INFO) {
		message := fmt.Sprintf(format, args...)

		logger.Info(message)

		loggingCallback(logging.INFO, message)
	}
}

func Notice(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.NOTICE) {
		message := fmt.Sprintf(format, args...)

		logger.Notice(message)

		loggingCallback(logging.NOTICE, message)
	}
}

func Warning(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.WARNING) {
		message := fmt.Sprintf(format, args...)

		logger.Warning(message)

		loggingCallback(logging.WARNING, message)
	}
}

func Error(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.ERROR) {
		message := fmt.Sprintf(format, args...)

		logger.Error(message)

		loggingCallback(logging.ERROR, message)
	}
}

func Critical(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.CRITICAL) {
		message := fmt.Sprintf(format, args...)

		logger.Critical(message)

		loggingCallback(logging.CRITICAL, message)
	}
}

func SetLevel(level logging.Level) {
	leveledBackend.SetLevel(level, "")
}

func init() {
	logsDirectory := filepath.Join(moonclient.Directory, "logs")

	err := os.MkdirAll(logsDirectory, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create the logs directory: %v", err.Error())
		os.Exit(1)
	}

	now := time.Now()
	logsFileName := fmt.Sprintf("MoonDeploy - %d-%d-%d %d-%d-%d %d.log",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond())

	logFilePath := filepath.Join(logsDirectory, logsFileName)

	targetOutput, err := os.Create(logFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create the log file! %v", err.Error())
		os.Exit(1)
	}

	backend := logging.NewLogBackend(targetOutput, "", 0)

	formatString := "%{time:15:04:05.000} - %{level} %{message}\n\n"

	format := logging.MustStringFormatter(formatString)

	backendFormatter := logging.NewBackendFormatter(backend, format)

	leveledBackend = logging.AddModuleLevel(backendFormatter)

	logger.SetBackend(leveledBackend)
}
