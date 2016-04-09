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

package apps

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/launchers"
	"github.com/giancosta86/moondeploy/v3/log"
)

func (app *App) CreateDesktopShortcut(launcher launchers.Launcher, referenceDescriptor descriptors.AppDescriptor) (err error) {
	desktopDir, err := caravel.GetUserDesktop()
	if err != nil {
		return err
	}

	if !caravel.DirectoryExists(desktopDir) {
		return fmt.Errorf("Expected desktop dir '%v' not found", desktopDir)
	}

	scriptFileName := fmt.Sprintf("%v.scpt", caravel.FormatFileName(referenceDescriptor.GetName()))
	log.Debug("Script file name: '%v'", scriptFileName)

	scriptFilePath := filepath.Join(desktopDir, scriptFileName)
	log.Debug("Script file to create: '%v'...", scriptFilePath)

	scriptGenerationCommand := exec.Command(
		"osacompile",
		"-e",
		fmt.Sprintf(`do shell script ""%v" "%v""`,
			launcher.GetExecutable(),
			app.GetLocalDescriptorPath()),
		"-o",
		scriptFilePath)

	log.Debug("Script command is: %v", scriptGenerationCommand)

	err = scriptGenerationCommand.Run()
	if err != nil {
		return err
	}

	log.Notice("Shortcut script created")

	return nil
}
