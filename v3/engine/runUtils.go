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
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3/apps"
	"github.com/giancosta86/moondeploy/v3/moonclient"
	"github.com/giancosta86/moondeploy/v3/custom"
	"github.com/giancosta86/moondeploy/v3/logging"
	"github.com/giancosta86/moondeploy/v3/ui"
)

func resolveAppDir(bootDescriptor apps.AppDescriptor, appGalleryDir string) (appDir string, err error) {
	baseURL := bootDescriptor.GetDeclaredBaseURL()

	hostComponent := strings.Replace(baseURL.Host, ":", "_", -1)

	appDirComponents := []string{
		appGalleryDir,
		hostComponent}

	trimmedBasePath := strings.Trim(baseURL.Path, "/")
	baseComponents := strings.Split(trimmedBasePath, "/")

	appDirComponents = append(appDirComponents, baseComponents...)

	if hostComponent == "github.com" &&
		len(appDirComponents) > 2 &&
		appDirComponents[len(appDirComponents)-2] == "releases" &&
		appDirComponents[len(appDirComponents)-1] == "latest" {
		appDirComponents = appDirComponents[0 : len(appDirComponents)-2]
	}

	appDir = filepath.Join(appDirComponents...)

	return appDir, nil
}

func ensureFirstRun(bootDescriptor apps.AppDescriptor, appDir string, userInterface ui.UserInterface) (err error) {
	var canRun bool
	if caravel.IsSecureURL(bootDescriptor.GetDeclaredBaseURL()) {
		canRun = userInterface.AskForSecureFirstRun(bootDescriptor)
	} else {
		canRun = userInterface.AskForUntrustedFirstRun(bootDescriptor)
	}

	if !canRun {
		return &ExecutionCanceled{}
	}

	logging.Notice("The user agreed to proceed")

	logging.Info("Ensuring the app dir is available...")
	err = os.MkdirAll(appDir, 0700)
	if err != nil {
		return err
	}
	logging.Notice("App dir available")

	return nil
}

func getLocalDescriptor(localDescriptorPath string) (localDescriptor apps.AppDescriptor) {
	if !caravel.FileExists(localDescriptorPath) {
		logging.Notice("The local descriptor is missing")
		return nil
	}

	logging.Notice("The local descriptor has been found! Deserializing...")
	localDescriptor, err := apps.NewAppDescriptorFromPath(localDescriptorPath)
	if err != nil {
		logging.Warning(err.Error())
		return nil
	}
	logging.Notice("Local descriptor deserialized")

	logging.Info("The local descriptor is: %#v", localDescriptor)

	return localDescriptor
}

func getRemoteDescriptor(bootDescriptor apps.AppDescriptor, localDescriptor apps.AppDescriptor, userInterface ui.UserInterface) (remoteDescriptor apps.AppDescriptor) {
	var remoteDescriptorURL *url.URL
	var err error

	if localDescriptor != nil {
		remoteDescriptorURL, err = localDescriptor.GetFileURL(localDescriptor.GetDescriptorFileName())
	} else {
		remoteDescriptorURL, err = bootDescriptor.GetFileURL(bootDescriptor.GetDescriptorFileName())
	}

	if err != nil {
		logging.Warning(err.Error())
		return nil
	}

	logging.Notice("The remote descriptor's URL is: %v", remoteDescriptorURL)

	logging.Info("Retrieving the remote descriptor...")
	remoteDescriptorBytes, err := caravel.RetrieveFromURL(remoteDescriptorURL)
	if err != nil {
		logging.Warning(err.Error())
		return nil
	}
	logging.Notice("Remote descriptor retrieved")

	logging.Info("Deserializing the remote descriptor...")
	remoteDescriptor, err = apps.NewAppDescriptorFromBytes(remoteDescriptorBytes)
	if err != nil {
		logging.Warning(err.Error())
		return nil
	}
	logging.Notice("Remote descriptor deserialized")

	logging.Notice("The remote descriptor is: %#v", remoteDescriptor)

	return remoteDescriptor
}

func chooseReferenceDescriptor(remoteDescriptor apps.AppDescriptor, localDescriptor apps.AppDescriptor) (referenceDescriptor apps.AppDescriptor, err error) {
	if remoteDescriptor == nil && localDescriptor == nil {
		return nil, fmt.Errorf("Cannot run the application: it is not installed and cannot be downloaded")
	}

	if remoteDescriptor == nil {
		if localDescriptor.IsSkipUpdateCheck() {
			logging.Info("The remote descriptor is missing as requested, so the local descriptor will be used")
		} else {
			logging.Warning("The remote descriptor is missing, so the local descriptor will be used")
		}
		return localDescriptor, nil
	}

	if localDescriptor == nil {
		logging.Notice("The local descriptor is missing, so the remote descriptor will be used")
		return remoteDescriptor, nil
	}

	if remoteDescriptor.GetAppVersion().NewerThan(localDescriptor.GetAppVersion()) {
		logging.Notice("Switching to the remote descriptor, as it is more recent")
		return remoteDescriptor, nil
	}

	logging.Notice("Keeping the local descriptor, as the remote descriptor is NOT more recent")
	return localDescriptor, nil
}

func prepareCommand(appDir string, appFilesDir string, commandLine []string) (command *exec.Cmd) {
	if caravel.DirectoryExists(appFilesDir) {
		os.Chdir(appFilesDir)
		logging.Notice("Files directory set as the current directory")
	} else {
		os.Chdir(appDir)
		logging.Notice("App directory set as the current directory")
	}

	logging.Info("Creating the command...")

	if len(commandLine) == 1 {
		return exec.Command(commandLine[0])
	}

	return exec.Command(commandLine[0], commandLine[1:]...)
}

func tryToSaveReferenceDescriptor(referenceDescriptor apps.AppDescriptor, localDescriptorPath string) (referenceDescriptorSaved bool) {
	logging.Info("Saving the reference descriptor as the local descriptor...")
	referenceDescriptorBytes, err := referenceDescriptor.GetBytes()
	if err != nil {
		logging.Error("Could not serialize the reference descriptor: %v", err)
		return false
	}

	err = ioutil.WriteFile(localDescriptorPath, referenceDescriptorBytes, 0600)
	if err != nil {
		logging.Error("Could not save the reference descriptor: %v", err)
		return false
	}

	logging.Notice("Reference descriptor saved")
	return true
}

func launchApp(command *exec.Cmd, settings *custom.Settings, userInterface ui.UserInterface) (err error) {
	logging.Info("Starting the app...")

	logging.Info("Hiding the user interface...")
	userInterface.HideLoader()
	logging.Notice("User interface hidden")

	if settings.SkipAppOutput {
		return command.Run()
	}
	var outputBytes []byte
	outputBytes, err = command.CombinedOutput()

	if outputBytes != nil && len(outputBytes) > 0 {
		fmt.Println("------------------------------")
		fmt.Printf("%s\n", outputBytes)
		fmt.Println("------------------------------")
	}

	return err
}

func getActualIconPath(referenceDescriptor apps.AppDescriptor, appFilesDir string) string {
	providedIconPath := referenceDescriptor.GetIconPath()

	if providedIconPath != "" {
		return filepath.Join(appFilesDir, providedIconPath)
	}

	return moonclient.GetIconPath()
}
