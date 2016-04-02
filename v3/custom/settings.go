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

package custom

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3/moonclient"
)

const userSettingsFileName = ".moondeploy.json"

const defaultGalleryDirName = "apps"
const defaultBufferSize = 512 * 1024
const defaultLoggingLevel = "INFO"
const defaultLoggingLevelValue = logging.INFO
const defaultSkipAppOutput = false

type Settings struct {
	GalleryDir    string
	BufferSize    int64
	LoggingLevel  string
	SkipAppOutput bool
}

func GetDefaultSettings() (settings *Settings, err error) {
	defaultGalleryDir, err := getDefaultGalleryDir()
	if err != nil {
		return nil, err
	}

	return &Settings{
		GalleryDir:    defaultGalleryDir,
		BufferSize:    defaultBufferSize,
		LoggingLevel:  defaultLoggingLevel,
		SkipAppOutput: defaultSkipAppOutput,
	}, nil
}

func getDefaultGalleryDir() (galleryDir string, err error) {
	galleryDir = filepath.Join(moonclient.Dir, defaultGalleryDirName)

	return galleryDir, nil
}

func ComputeUserSettings() (userSettings *Settings, err error) {
	userDir, err := caravel.GetUserDirectory()
	if err != nil {
		return GetDefaultSettings()
	}

	userSettingsPath := filepath.Join(userDir, userSettingsFileName)
	if !caravel.FileExists(userSettingsPath) {
		return GetDefaultSettings()
	}

	userSettingsBytes, err := ioutil.ReadFile(userSettingsPath)
	if err != nil {
		return GetDefaultSettings()
	}

	return deserializeAndMergeSettings(userSettingsBytes)
}

func deserializeAndMergeSettings(settingsBytes []byte) (settings *Settings, err error) {
	settings = &Settings{}

	if settingsBytes != nil {
		json.Unmarshal(settingsBytes, settings)
	}

	defaultSettings, err := GetDefaultSettings()

	if settings.GalleryDir == "" {
		if err != nil {
			return nil, err
		}

		settings.GalleryDir = defaultSettings.GalleryDir
	}

	if settings.BufferSize == 0 {
		if err != nil {
			return nil, err
		}

		settings.BufferSize = defaultSettings.BufferSize
	}

	if settings.LoggingLevel == "" {
		if err != nil {
			return nil, err
		}

		settings.LoggingLevel = defaultSettings.LoggingLevel
	}

	//No need to keep track of bool values

	return settings, nil
}

func (settings *Settings) GetLoggingLevel() (level logging.Level) {
	lowercaseLevelString := strings.ToLower(settings.LoggingLevel)

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
		return defaultLoggingLevelValue
	}
}
