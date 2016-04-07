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
	"io"
	"os"

	"github.com/op/go-logging"
)

type LogCallback func(level logging.Level, message string)

var logger logging.Logger
var leveledBackend logging.LeveledBackend

var logCallback LogCallback = func(level logging.Level, message string) {}

func SetCallback(callback LogCallback) {
	logCallback = callback
}

func Debug(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.DEBUG) {
		message := fmt.Sprintf(format, args...)

		logger.Debug(message)

		logCallback(logging.DEBUG, message)
	}
}

func Info(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.INFO) {
		message := fmt.Sprintf(format, args...)

		logger.Info(message)

		logCallback(logging.INFO, message)
	}
}

func Notice(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.NOTICE) {
		message := fmt.Sprintf(format, args...)

		logger.Notice(message)

		logCallback(logging.NOTICE, message)
	}
}

func Warning(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.WARNING) {
		message := fmt.Sprintf(format, args...)

		logger.Warning(message)

		logCallback(logging.WARNING, message)
	}
}

func Error(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.ERROR) {
		message := fmt.Sprintf(format, args...)

		logger.Error(message)

		logCallback(logging.ERROR, message)
	}
}

func Critical(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.CRITICAL) {
		message := fmt.Sprintf(format, args...)

		logger.Critical(message)

		logCallback(logging.CRITICAL, message)
	}
}

func SetLevel(level logging.Level) {
	leveledBackend.SetLevel(level, "")
}

func Setup(logWriter io.Writer) {
	backend := logging.NewLogBackend(logWriter, "", 0)

	formatString := "%{time:15:04:05.000} %{level} - %{message}\n\n"

	formatter := logging.MustStringFormatter(formatString)

	backendFormatter := logging.NewBackendFormatter(backend, formatter)

	leveledBackend = logging.AddModuleLevel(backendFormatter)

	logger.SetBackend(leveledBackend)
}

func init() {
	Setup(os.Stdout)
}
