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
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/launchers"
	"github.com/giancosta86/moondeploy/v3/log"
)

const windowsShortcutContent = `
	set WshShell = WScript.CreateObject("WScript.Shell")
	set shellLink = WshShell.CreateShortcut("%v")
	shellLink.TargetPath = "%v"
	shellLink.Description = "%v"
	shellLink.IconLocation = "%v"
	shellLink.WorkingDirectory = "%v"
	shellLink.Save`

func (app *App) CreateDesktopShortcut(launcher *launchers.Launcher, referenceDescriptor descriptors.AppDescriptor) (err error) {
	desktopDir, err := caravel.GetUserDesktop()
	if err != nil {
		return err
	}

	shortcutName := caravel.FormatFileName(referenceDescriptor.GetName()) + ".lnk"
	log.Info("Shortcut name: '%v'", shortcutName)

	shortcutPath := filepath.Join(desktopDir, shortcutName)
	log.Info("Shortcut path: '%v'", shortcutPath)

	log.Info("Creating temp file for script...")

	salt := rand.Int63()
	tempFileName := fmt.Sprintf("shortcutScript_%v_%v.vbs", time.Now().Unix(), salt)
	tempFilePath := filepath.Join(os.TempDir(), tempFileName)
	tempFile, err := os.Create(tempFilePath)

	if err != nil {
		return err
	}
	defer func() {
		tempFile.Close()

		tempRemovalErr := os.Remove(tempFilePath)
		if tempRemovalErr != nil {
			log.Warning("Cannot remove the temp script: %v", tempFilePath)
		} else {
			log.Info("Temp script successfully removed")
		}

		if err != nil {
			os.Remove(shortcutPath)
		}
	}()

	log.Info("Temp script file created: %v", tempFilePath)
	actualIconPath := app.GetActualIconPath(launcher)
	log.Info("Actual icon path: '%v'", actualIconPath)

	workingDirectory := filepath.Dir(app.localDescriptorPath)
	log.Info("Working directory: '%v'", workingDirectory)

	shortcutScript := fmt.Sprintf(windowsShortcutContent,
		shortcutPath,
		app.localDescriptorPath,
		referenceDescriptor.GetDescription(),
		actualIconPath,
		workingDirectory)

	log.Info("Writing script temp file...")
	tempFile.Write([]byte(shortcutScript))
	tempFile.Close()
	log.Info("Temp script ready")

	log.Info("Now executing the temp script...")

	command := exec.Command("wscript", "/b", "/nologo", tempFilePath)

	err = command.Run()
	if err != nil {
		return err
	}

	if !command.ProcessState.Success() {
		return fmt.Errorf("The script did not run successfully")
	}

	log.Notice("The script was successful")

	return nil
}
