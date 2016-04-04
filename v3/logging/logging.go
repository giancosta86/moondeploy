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

package logging

import (
	"fmt"
	"os"
	"runtime"

	"github.com/op/go-logging"
)

var logger logging.Logger
var leveledBackend logging.LeveledBackend

type LoggingCallback func(message string)

var loggingCallback LoggingCallback = func(message string) {}

var outputEnabled = true

func SetOutputEnabled(newValue bool) {
	outputEnabled = newValue
}

func SetCallback(callback LoggingCallback) {
	loggingCallback = callback
}

func Debug(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.DEBUG) {
		message := fmt.Sprintf(format, args...)

		if outputEnabled {
			logger.Debug(message)
		}

		loggingCallback(message)
	}
}

func Info(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.INFO) {
		message := fmt.Sprintf(format, args...)

		if outputEnabled {
			logger.Info(message)
		}

		loggingCallback(message)
	}
}

func Notice(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.NOTICE) {
		message := fmt.Sprintf(format, args...)

		if outputEnabled {
			logger.Notice(message)
		}

		loggingCallback(message)
	}
}

func Warning(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.WARNING) {
		message := fmt.Sprintf(format, args...)

		if outputEnabled {
			logger.Warning(message)
		}

		loggingCallback(message)
	}
}

func Error(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.ERROR) {
		message := fmt.Sprintf(format, args...)

		if outputEnabled {
			logger.Error(message)
		}

		loggingCallback(message)
	}
}

func Critical(format string, args ...interface{}) {
	if logger.IsEnabledFor(logging.CRITICAL) {
		message := fmt.Sprintf(format, args...)

		if outputEnabled {
			logger.Critical(message)
		}

		loggingCallback(message)
	}
}

func SetLevel(level logging.Level) {
	leveledBackend.SetLevel(level, "")
}

func init() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)

	var formatString string

	if runtime.GOOS != "windows" {
		formatString = "%{color}%{time:15:04:05.000} ▶ %{level}%{color:reset} %{message}"
	} else {
		formatString = "%{time:15:04:05.000} - %{level} %{message}"
	}

	format := logging.MustStringFormatter(formatString)

	backendFormatter := logging.NewBackendFormatter(backend, format)

	leveledBackend = logging.AddModuleLevel(backendFormatter)

	logger.SetBackend(leveledBackend)
}
