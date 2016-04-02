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
	"io/ioutil"
	"os"

	"github.com/giancosta86/caravel"
	"github.com/giancosta86/moondeploy/apps"
	"github.com/giancosta86/moondeploy/custom"
	"github.com/giancosta86/moondeploy/logging"
	"github.com/giancosta86/moondeploy/ui"
)

func checkAppFiles(
	remoteDescriptor apps.AppDescriptor,
	localDescriptorPath string,
	localDescriptor apps.AppDescriptor,
	appFilesDir string,
	settings *custom.Settings,
	userInterface ui.UserInterface) (err error) {

	packagesToUpdate := getPackagesToUpdate(remoteDescriptor, localDescriptor)

	if len(packagesToUpdate) == 0 {
		logging.Notice("All the packages are up-to-date")
		return nil
	}

	if caravel.FileExists(localDescriptorPath) {
		logging.Info("Deleting the local descriptor before starting the update process...")
		err = os.Remove(localDescriptorPath)
		if err != nil {
			return err
		}
		logging.Notice("Local descriptor deleted")
	}

	retrieveAllPackages := (len(packagesToUpdate) == len(remoteDescriptor.GetPackageVersions()))
	logging.Notice("Must retrieve all the remote packages? %v", retrieveAllPackages)

	if retrieveAllPackages {
		logging.Info("Removing app files dir...")
		err = os.RemoveAll(appFilesDir)
		if err != nil {
			return err
		}
		logging.Notice("App files dir removed")
	}

	for packageIndex, packageName := range packagesToUpdate {
		userInterface.SetHeader(
			fmt.Sprintf("Updating package %v of %v: %v",
				packageIndex+1,
				len(packagesToUpdate),
				packageName))

		err = installPackage(
			remoteDescriptor,
			packageName,
			appFilesDir,
			settings,
			func(retrievedSize int64, totalSize int64) {
				userInterface.SetProgress(float64(retrievedSize) / float64(totalSize))
				logging.Info("Retrieved: %v / %v bytes", retrievedSize, totalSize)
			})
		if err != nil {
			return err
		}
	}

	return nil
}

func getPackagesToUpdate(remoteDescriptor apps.AppDescriptor, localDescriptor apps.AppDescriptor) []string {
	if localDescriptor == nil {
		packagesToUpdate := []string{}

		for packageName := range remoteDescriptor.GetPackageVersions() {
			packagesToUpdate = append(packagesToUpdate, packageName)
		}

		return packagesToUpdate
	}

	if !remoteDescriptor.GetAppVersion().NewerThan(localDescriptor.GetAppVersion()) {
		return []string{}
	}

	packagesToUpdate := []string{}

	for remotePackageName, remotePackageVersion := range remoteDescriptor.GetPackageVersions() {
		localPackageVersion := localDescriptor.GetPackageVersions()[remotePackageName]

		if remotePackageVersion == nil ||
			localPackageVersion == nil ||
			remotePackageVersion.NewerThan(localPackageVersion) {
			packagesToUpdate = append(packagesToUpdate, remotePackageName)
		}
	}

	return packagesToUpdate
}

func installPackage(
	remoteDescriptor apps.AppDescriptor,
	packageName string,
	appFilesDir string,
	settings *custom.Settings,
	progressCallback caravel.RetrievalProgressCallback) (err error) {

	packageURL, err := remoteDescriptor.GetFileURL(packageName)
	if err != nil {
		return err
	}

	logging.Info("Creating package temp file...")
	packageTempFile, err := ioutil.TempFile(os.TempDir(), packageName)
	if err != nil {
		return err
	}
	packageTempFilePath := packageTempFile.Name()
	logging.Info("Package temp file created '%v'", packageTempFilePath)

	defer func() {
		packageTempFile.Close()

		logging.Info("Deleting package temp file: '%v'", packageTempFilePath)
		tempFileRemovalErr := os.Remove(packageTempFilePath)
		if tempFileRemovalErr != nil {
			logging.Warning("Could not remove the package temp file! '%v'", tempFileRemovalErr)
		} else {
			logging.Notice("Package temp file removed")
		}
	}()

	logging.Info("Retrieving package: %v", packageURL)
	err = caravel.RetrieveChunksFromURL(packageURL, packageTempFile, settings.BufferSize, progressCallback)
	if err != nil {
		return err
	}
	logging.Notice("Package retrieved")

	logging.Info("Closing the package temp file...")
	packageTempFile.Close()
	if err != nil {
		return err
	}
	logging.Notice("Package temp file closed")

	err = os.MkdirAll(appFilesDir, 0700)
	if err != nil {
		return err
	}

	logging.Info("Extracting the package. Skipping levels: %v...", remoteDescriptor.GetSkipPackageLevels())
	err = caravel.ExtractZipSkipLevels(packageTempFilePath, appFilesDir, remoteDescriptor.GetSkipPackageLevels())
	if err != nil {
		return err
	}
	logging.Notice("Package extracted")

	return nil
}
