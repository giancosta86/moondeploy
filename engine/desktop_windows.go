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

package engine

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/apps"
	"github.com/giancosta86/moondeploy/logging"
)

const windowsShortcutContent = `
	set WshShell = WScript.CreateObject("WScript.Shell")
	set shellLink = WshShell.CreateShortcut("%v")
	shellLink.TargetPath = "%v"
	shellLink.Description = "%v"
	shellLink.IconLocation = "%v"
	shellLink.WorkingDirectory = "%v"
	shellLink.Save`

func createDesktopShortcut(appFilesDir string, localDescriptorPath string, referenceDescriptor *apps.AppDescriptor) (err error) {
	desktopDir, err := caravel.GetUserDesktop()
	if err != nil {
		return err
	}

	shortcutName := caravel.FormatFileName(referenceDescriptor.Name) + ".lnk"
	logging.Info("Shortcut name: '%v'", shortcutName)

	shortcutPath := filepath.Join(desktopDir, shortcutName)
	logging.Info("Shortcut path: '%v'", shortcutPath)

	logging.Info("Creating temp file for script...")

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
			logging.Warning("Cannot remove the temp script: %v", tempFilePath)
		} else {
			logging.Info("Temp script successfully removed")
		}

		if err != nil {
			os.Remove(shortcutPath)
		}
	}()

	logging.Info("Temp script file created: %v", tempFilePath)
	actualIconPath := referenceDescriptor.GetActualIconPath(appFilesDir)
	logging.Info("Actual icon path: '%v'", actualIconPath)

	workingDirectory := filepath.Dir(localDescriptorPath)
	logging.Info("Working directory: '%v'", workingDirectory)

	shortcutScript := fmt.Sprintf(windowsShortcutContent,
		shortcutPath,
		localDescriptorPath,
		referenceDescriptor.Description,
		actualIconPath,
		workingDirectory)

	logging.Info("Writing script temp file...")
	tempFile.Write([]byte(shortcutScript))
	tempFile.Close()
	logging.Info("Temp script ready")

	logging.Info("Now executing the temp script...")

	command := exec.Command("wscript", "/b", "/nologo", tempFilePath)

	err = command.Run()
	if err != nil {
		return err
	}

	if !command.ProcessState.Success() {
		return fmt.Errorf("The script did not run successfully")
	}

	logging.Notice("The script was successful")

	return nil
}
