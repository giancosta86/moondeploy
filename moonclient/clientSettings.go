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
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3"
	"github.com/giancosta86/moondeploy/v3/log"
)

const userSettingsFileName = ".moondeploy.json"

const defaultLocalDirName = "MoonDeploy"
const galleryDirName = "apps"
const logsDirName = "logs"

const defaultBufferSize = 1024 * 1024
const defaultSkipAppOutput = false

const defaultLoggingLevel = logging.DEBUG

const defaultBackgroundColor = 195
const defaultForegroundColor = 22

const defaultLogMaxAgeInHours = 120

type rawMoonSettingsStruct struct {
	LocalDirectory   string
	BufferSize       int64
	LoggingLevel     string
	SkipAppOutput    bool
	BackgroundColor  int
	ForegroundColor  int
	LogMaxAgeInHours int
}

type MoonSettings struct {
	localDirectory   string
	galleryDirectory string
	logsDirectory    string
	bufferSize       int64
	loggingLevel     logging.Level
	skipAppOutput    bool
	backgroundColor  int
	foregroundColor  int
	logMaxAgeInHours int
}

var moonSettings *MoonSettings

func (settings *MoonSettings) GetLocalDirectory() string {
	return settings.localDirectory
}

func (settings *MoonSettings) GetGalleryDirectory() string {
	return settings.galleryDirectory
}

func (settings *MoonSettings) GetLogsDirectory() string {
	return settings.logsDirectory
}

func (settings *MoonSettings) GetBufferSize() int64 {
	return settings.bufferSize
}

func (settings *MoonSettings) GetLoggingLevel() logging.Level {
	return settings.loggingLevel
}

func (settings *MoonSettings) IsSkipAppOutput() bool {
	return settings.skipAppOutput
}

func (settings *MoonSettings) GetBackgroundColor() int {
	return settings.backgroundColor
}

func (settings *MoonSettings) GetForegroundColor() int {
	return settings.foregroundColor
}

func (settings *MoonSettings) GetLogMaxAgeInHours() int {
	return settings.logMaxAgeInHours
}

func getRawMoonSettings() (rawMoonSettings *rawMoonSettingsStruct) {
	rawMoonSettings = &rawMoonSettingsStruct{
		BackgroundColor:  -1,
		ForegroundColor:  -1,
		LogMaxAgeInHours: defaultLogMaxAgeInHours,
	}

	userDir, err := caravel.GetUserDirectory()
	if err != nil {
		return rawMoonSettings
	}

	userSettingsPath := filepath.Join(userDir, userSettingsFileName)
	if !caravel.FileExists(userSettingsPath) {
		return rawMoonSettings
	}

	rawSettingsBytes, err := ioutil.ReadFile(userSettingsPath)
	if err != nil {
		return rawMoonSettings
	}

	err = json.Unmarshal(rawSettingsBytes, rawMoonSettings)
	if err != nil {
		return &rawMoonSettingsStruct{}
	}

	log.Debug("Settings file content: %#v", rawMoonSettings)

	return rawMoonSettings
}

func parseLoggingLevel(loggingLevelName string) (level logging.Level) {
	lowercaseLevelString := strings.ToLower(loggingLevelName)

	switch lowercaseLevelString {
	case "debug":
		return logging.DEBUG
	case "info":
		return logging.INFO
	case "notice":
		return logging.NOTICE
	case "warning":
		return logging.WARNING
	case "error":
		return logging.ERROR
	case "critical":
		return logging.CRITICAL

	default:
		return defaultLoggingLevel
	}
}

func getMoonSettings() *MoonSettings {
	if moonSettings != nil {
		return moonSettings
	}

	rawMoonSettings := getRawMoonSettings()

	moonSettings = &MoonSettings{}

	if rawMoonSettings.LocalDirectory != "" {
		moonSettings.localDirectory = rawMoonSettings.LocalDirectory
	} else {
		userDir, err := caravel.GetUserDirectory()
		if err != nil {
			log.Error("Cannot retrieve the user's directory")
			os.Exit(v3.ExitCodeError)
		}

		moonSettings.localDirectory = filepath.Join(userDir, defaultLocalDirName)
	}

	moonSettings.galleryDirectory = filepath.Join(moonSettings.localDirectory, galleryDirName)
	moonSettings.logsDirectory = filepath.Join(moonSettings.localDirectory, logsDirName)

	if rawMoonSettings.BufferSize > 0 {
		moonSettings.bufferSize = rawMoonSettings.BufferSize
	} else {
		moonSettings.bufferSize = defaultBufferSize
	}

	moonSettings.loggingLevel = parseLoggingLevel(rawMoonSettings.LoggingLevel)

	moonSettings.skipAppOutput = rawMoonSettings.SkipAppOutput

	if 0 <= rawMoonSettings.BackgroundColor && rawMoonSettings.BackgroundColor <= 255 {
		moonSettings.backgroundColor = rawMoonSettings.BackgroundColor
	} else {
		moonSettings.backgroundColor = defaultBackgroundColor
	}

	if 0 <= rawMoonSettings.ForegroundColor && rawMoonSettings.ForegroundColor <= 255 {
		moonSettings.foregroundColor = rawMoonSettings.ForegroundColor
	} else {
		moonSettings.foregroundColor = defaultForegroundColor
	}

	return moonSettings
}
