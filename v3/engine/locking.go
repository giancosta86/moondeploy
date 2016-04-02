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
	"os"
	"path/filepath"

	"github.com/giancosta86/LockAPI/lockapi"
	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3/apps"
	"github.com/giancosta86/moondeploy/v3/logging"
)

func lockAppDir(appDir string) (lockFile *os.File, err error) {
	lockFilePath := filepath.Join(appDir, apps.LockFileName)

	logging.Info("The lock file is: %v", lockFilePath)

	logging.Info("Opening the lock file...")
	lockFile, err = os.OpenFile(lockFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	logging.Info("Obtaining the API lock...")
	err = lockapi.TryLockFile(lockFile)
	if err != nil {
		lockFile.Close()
		return nil, err
	}
	logging.Notice("Lock acquired")

	return lockFile, nil
}

func unlockAppDir(lockFile *os.File) (err error) {
	if !caravel.FileExists(lockFile.Name()) {
		return nil
	}

	logging.Info("Releasing the API lock...")
	err = lockapi.UnlockFile(lockFile)
	if err != nil {
		return err
	}
	logging.Notice("Lock released")

	logging.Info("Closing lock file...")
	err = lockFile.Close()
	if err != nil {
		return err
	}
	logging.Notice("Lock file closed")

	logging.Info("Deleting lock file...")
	err = os.Remove(lockFile.Name())
	if err != nil {
		return err
	}
	logging.Notice("Lock file deleted")

	return nil
}
