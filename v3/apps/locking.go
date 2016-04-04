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
	"os"
	"path/filepath"

	"github.com/giancosta86/LockAPI/lockapi"
	"github.com/giancosta86/caravel"
	"github.com/giancosta86/moondeploy/v3/log"
)

func (app *App) LockDirectory() (err error) {
	lockFilePath := filepath.Join(app.Directory, lockFileName)

	log.Info("The lock file is: %v", lockFilePath)

	log.Info("Opening the lock file...")
	lockFile, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	log.Info("Obtaining the API lock...")
	err = lockapi.TryLockFile(lockFile)
	if err != nil {
		lockFile.Close()
		return err
	}
	log.Notice("Lock acquired")

	app.lockFile = lockFile

	return nil
}

func (app *App) UnlockDirectory() (err error) {
	if app.lockFile == nil {
		return nil
	}

	if !caravel.FileExists(app.lockFile.Name()) {
		return nil
	}

	log.Info("Releasing the API lock...")
	err = lockapi.UnlockFile(app.lockFile)
	if err != nil {
		return err
	}
	log.Notice("Lock released")

	log.Info("Closing lock file...")
	err = app.lockFile.Close()
	if err != nil {
		return err
	}
	log.Notice("Lock file closed")

	log.Info("Deleting lock file...")
	err = os.Remove(app.lockFile.Name())
	if err != nil {
		return err
	}
	log.Notice("Lock file deleted")

	app.lockFile = nil

	return nil
}
