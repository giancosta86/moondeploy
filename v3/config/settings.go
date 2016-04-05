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

package config

import (
	"strings"

	"github.com/op/go-logging"
)

const DefaultLoggingLevel = "INFO"
const DefaultLoggingLevelValue = logging.INFO

type Settings struct {
	GalleryDir    string
	BufferSize    int64
	LoggingLevel  string
	SkipAppOutput bool
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
		return DefaultLoggingLevelValue
	}
}
