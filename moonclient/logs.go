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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/giancosta86/caravel"
	"github.com/giancosta86/moondeploy/v3"
	"github.com/giancosta86/moondeploy/v3/config"
	"github.com/giancosta86/moondeploy/v3/log"
)

func initializeLogging(settings config.Settings) {
	log.SetLevel(settings.GetLoggingLevel())

	logsDirectory := settings.GetLogsDirectory()
	log.Debug("Logs directory is: '%v'", logsDirectory)

	tryToRemoveLogs(logsDirectory, settings)
	ensureLogsDirectory(logsDirectory)

	logFile := createLogFile(logsDirectory)

	log.Debug("Now redirecting log lines to file: '%v'", logFile.Name())
	log.Setup(logFile)
}

func tryToRemoveLogs(logsDirectory string, settings config.Settings) {
	if !caravel.DirectoryExists(logsDirectory) {
		return
	}

	logFiles, err := ioutil.ReadDir(logsDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot list the logs directory: %v\n", err)
		return
	}

	now := time.Now()

	for _, logFile := range logFiles {
		logFileAge := now.Sub(logFile.ModTime())

		if logFileAge.Hours() > float64(settings.GetLogMaxAgeInHours()) {
			logFilePath := filepath.Join(logsDirectory, logFile.Name())
			err = os.Remove(logFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot delete log file: '%v'\n", logFile.Name())
			}
		}
	}
}

func ensureLogsDirectory(logsDirectory string) {
	err := os.MkdirAll(logsDirectory, 0700)
	if err != nil {
		log.Error("Cannot create the logs directory ('%v'). %v", logsDirectory, err)
		os.Exit(v3.ExitCodeError)
	}
}

func createLogFile(logsDirectory string) *os.File {
	now := time.Now()

	logFileName := fmt.Sprintf("MoonDeploy - %d-%d-%d %d-%d-%d %d.log",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond())

	logFilePath := filepath.Join(logsDirectory, logFileName)

	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Error("Cannot create the log file '%v'! %v", logFilePath, err.Error())
		os.Exit(v3.ExitCodeError)
	}

	return logFile
}
