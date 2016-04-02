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
	"os"
	"path/filepath"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy"
	"github.com/giancosta86/moondeploy/apps"
	"github.com/giancosta86/moondeploy/logging"
)

const linuxShortcutContent = `
[Desktop Entry]
Encoding=UTF-8
Name=%v
Comment=%v
Exec="%v" "%v"
Icon=%v
Version=1.0
Type=Application
Terminal=0
`

func createDesktopShortcut(appFilesDir string, localDescriptorPath string, referenceDescriptor apps.AppDescriptor) (err error) {
	desktopDir, err := caravel.GetUserDesktop()
	if err != nil {
		return err
	}

	if !caravel.DirectoryExists(desktopDir) {
		return fmt.Errorf("Expected desktop dir '%v' not found", desktopDir)
	}

	shortcutFileName := caravel.FormatFileName(referenceDescriptor.GetName()) + ".desktop"
	logging.Info("Shortcut name: '%v'", shortcutFileName)

	shortcutFilePath := filepath.Join(desktopDir, shortcutFileName)

	logging.Info("Creating desktop shortcut: '%v'...", shortcutFilePath)

	shortcutFile, err := os.OpenFile(shortcutFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer func() {
		shortcutFile.Close()
		if err != nil {
			os.Remove(shortcutFilePath)
		}
	}()

	actualIconPath := getActualIconPath(referenceDescriptor, appFilesDir)

	shortcutContent := fmt.Sprintf(linuxShortcutContent,
		referenceDescriptor.GetName(),
		referenceDescriptor.GetDescription(),
		moondeploy.Executable,
		localDescriptorPath,
		actualIconPath)

	_, err = shortcutFile.Write([]byte(shortcutContent))
	if err != nil {
		return err
	}

	logging.Notice("Desktop shortcut created")

	return nil
}
