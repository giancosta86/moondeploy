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

package moonclient

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/giancosta86/caravel"
	"github.com/giancosta86/moondeploy/v3/config"
)

const userSettingsFileName = ".moondeploy.json"

const defaultGalleryDirName = "apps"
const defaultBufferSize = 512 * 1024
const defaultSkipAppOutput = false

func getDefaultGalleryDir() string {
	return filepath.Join(Directory, defaultGalleryDirName)
}

func GetDefaultSettings() (defaultSettings *config.Settings, err error) {
	defaultGalleryDir := getDefaultGalleryDir()

	return &config.Settings{
		GalleryDir:    defaultGalleryDir,
		BufferSize:    defaultBufferSize,
		LoggingLevel:  config.DefaultLoggingLevel,
		SkipAppOutput: defaultSkipAppOutput,
	}, nil
}

func ComputeUserSettings() (userSettings *config.Settings, err error) {
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

func deserializeAndMergeSettings(settingsBytes []byte) (settings *config.Settings, err error) {
	defaultSettings, err := GetDefaultSettings()

	settings = &config.Settings{}

	if settingsBytes != nil {
		json.Unmarshal(settingsBytes, settings)
	}

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
