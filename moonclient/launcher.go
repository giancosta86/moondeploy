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
	"path/filepath"
	"runtime"

	"github.com/kardianos/osext"

	"github.com/giancosta86/moondeploy"
	"github.com/giancosta86/moondeploy/v3/config"
)

type MoonLauncher struct {
	name          string
	title         string
	executable    string
	directory     string
	iconPathAsIco string
	iconPathAsPng string
	settings      config.Settings
}

var moonLauncher *MoonLauncher

func (launcher *MoonLauncher) GetName() string {
	return launcher.name
}

func (launcher *MoonLauncher) GetTitle() string {
	return launcher.title
}

func (launcher *MoonLauncher) GetExecutable() string {
	return launcher.executable
}

func (launcher *MoonLauncher) GetDirectory() string {
	return launcher.directory
}

func (launcher *MoonLauncher) GetIconPath() string {
	if runtime.GOOS == "windows" {
		return launcher.GetIconPathAsIco()
	}

	return launcher.GetIconPathAsPng()
}

func (launcher *MoonLauncher) GetIconPathAsIco() string {
	return launcher.iconPathAsIco
}

func (launcher *MoonLauncher) GetIconPathAsPng() string {
	return launcher.iconPathAsPng
}

func (launcher *MoonLauncher) GetSettings() config.Settings {
	return launcher.settings
}

func getMoonLauncher() *MoonLauncher {
	if moonLauncher != nil {
		return moonLauncher
	}

	moonLauncher = &MoonLauncher{
		name: "MoonDeploy",
	}

	moonLauncher.title = fmt.Sprintf("%v %v", moonLauncher.name, moondeploy.Version)

	var err error
	moonLauncher.executable, err = osext.Executable()
	if err != nil {
		panic(err)
	}

	moonLauncher.directory, err = osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

	moonLauncher.iconPathAsIco = filepath.Join(moonLauncher.directory, "moondeploy.ico")
	moonLauncher.iconPathAsPng = filepath.Join(moonLauncher.directory, "moondeploy.png")

	moonLauncher.settings = getMoonSettings()

	return moonLauncher
}
