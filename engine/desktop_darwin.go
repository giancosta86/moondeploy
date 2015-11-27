/*§
  ===========================================================================
  MoonDeploy
  ===========================================================================
  Copyright (C) 2015 Gianluca Costa
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

const macScriptContent = `#!/bin/bash
"%v" "%v"
`

func createDesktopShortcut(appFilesDir string, localDescriptorPath string, referenceDescriptor *apps.AppDescriptor) (err error) {
	desktopDir, err := caravel.GetUserDesktop()
	if err != nil {
		return err
	}

	if !caravel.DirectoryExists(desktopDir) {
		return fmt.Errorf("Expected desktop dir '%v' not found", desktopDir)
	}

	scriptFileName := caravel.FormatFileName(referenceDescriptor.Name)
	logging.Info("BASH shortcut name: '%v'", scriptFileName)

	scriptFilePath := filepath.Join(desktopDir, scriptFileName)
	logging.Info("Creating BASH shortcut: '%v'...", scriptFilePath)

	scriptFile, err := os.OpenFile(scriptFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer func() {
		scriptFile.Close()

		if err != nil {
			os.Remove(scriptFilePath)
		}
	}()

	scriptContent := fmt.Sprintf(macScriptContent,
		moondeploy.Executable,
		localDescriptorPath)

	_, err = scriptFile.Write([]byte(scriptContent))
	if err != nil {
		return err
	}

	logging.Notice("BASH shortcut script created")

	return nil
}
