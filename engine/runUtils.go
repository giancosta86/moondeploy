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

	"github.com/giancosta86/moondeploy/apps"
	"github.com/giancosta86/moondeploy/custom"
	"github.com/giancosta86/moondeploy/gitHubUtils"
	"github.com/giancosta86/moondeploy/logging"
	"github.com/giancosta86/moondeploy/ui"
)

func resolveAppDir(bootDescriptor *apps.AppDescriptor, appGalleryDir string) (appDir string, err error) {
	hostComponent := strings.Replace(bootDescriptor.BaseURL.Host, ":", "_", -1)

	appDirComponents := []string{
		appGalleryDir,
		hostComponent}

	trimmedBasePath := strings.Trim(bootDescriptor.BaseURL.Path, "/")
	baseComponents := strings.Split(trimmedBasePath, "/")

	appDirComponents = append(appDirComponents, baseComponents...)

	appDir = filepath.Join(appDirComponents...)

	return appDir, nil
}

func ensureFirstRun(bootDescriptor *apps.AppDescriptor, appDir string, userInterface ui.UserInterface) (err error) {
	var canRun bool
	if caravel.IsSecureURL(bootDescriptor.BaseURL) {
		canRun = userInterface.AskForSecureFirstRun(bootDescriptor)
	} else {
		canRun = userInterface.AskForUntrustedFirstRun(bootDescriptor)
	}

	if !canRun {
		return &ExecutionCanceled{}
	}

	logging.Notice("The user agreed")

	logging.Info("Ensuring the app dir is available...")
	err = os.MkdirAll(appDir, 0700)
	if err != nil {
		return err
	}
	logging.Notice("App dir available")

	return nil
}

func getLocalDescriptor(localDescriptorPath string) (localDescriptor *apps.AppDescriptor) {
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

	logging.Info("Validating local descriptor...")
	err = localDescriptor.Validate()
	if err != nil {
		logging.Warning(err.Error())
		return nil
	}
	logging.Notice("Local descriptor valid")

	return localDescriptor
}

func getRemoteDescriptor(bootDescriptor *apps.AppDescriptor, localDescriptor *apps.AppDescriptor, userInterface ui.UserInterface) (remoteDescriptor *apps.AppDescriptor) {
	var remoteDescriptorURL *url.URL
	var err error

	logging.Info("Checking if the Base URL points to the *latest* release of a GitHub repo...")
	gitHubLatestRemoteDescriptorInfo := gitHubUtils.GetLatestRemoteDescriptorInfo(bootDescriptor.BaseURL)
	if gitHubLatestRemoteDescriptorInfo != nil {
		logging.Notice("The given base URL actually references version '%v', whose descriptor is at URL: '%v'",
			gitHubLatestRemoteDescriptorInfo.Version,
			gitHubLatestRemoteDescriptorInfo.DescriptorURL)

		if localDescriptor != nil && !gitHubLatestRemoteDescriptorInfo.Version.NewerThan(localDescriptor.Version) {
			logging.Notice("The remote descriptor is not newer than the local descriptor")
			return nil
		}

		remoteDescriptorURL = gitHubLatestRemoteDescriptorInfo.DescriptorURL
		logging.Notice("The remote descriptor will be downloaded from the new URL: '%v'", remoteDescriptorURL)

	} else {
		logging.Notice("The remote descriptor is NOT hosted on a GitHub *latest* release")

		remoteDescriptorURL, err = bootDescriptor.GetBaseFileURL(apps.DescriptorFileName)
		if err != nil {
			logging.Warning(err.Error())
			return nil
		}
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

	if gitHubLatestRemoteDescriptorInfo != nil {
		if remoteDescriptor.Version == nil || gitHubLatestRemoteDescriptorInfo.Version.CompareTo(remoteDescriptor.Version) != 0 {
			logging.Warning("The latest version returned by GitHub (%v) and the remote descriptor version (%v) do not match",
				gitHubLatestRemoteDescriptorInfo.Version,
				remoteDescriptor.Version)

			return nil
		}

		remoteDescriptorPathComponents := strings.Split(
			gitHubLatestRemoteDescriptorInfo.DescriptorURL.Path,
			"/")
		newBaseURLPathComponents := remoteDescriptorPathComponents[0 : len(remoteDescriptorPathComponents)-1]
		newBaseURLPath := strings.Join(newBaseURLPathComponents, "/") + "/"

		newBaseURLPathAsURL, err := url.Parse(newBaseURLPath)
		if err != nil {
			logging.Warning(err.Error())
			return nil
		}

		newBaseURL := remoteDescriptorURL.ResolveReference(newBaseURLPathAsURL)

		logging.Notice("The new base URL is: %v", newBaseURL)

		bootDescriptor.BaseURL = newBaseURL
		remoteDescriptor.BaseURL = newBaseURL

		if localDescriptor != nil {
			localDescriptor.BaseURL = newBaseURL
		}
	}

	logging.Notice("The remote descriptor is: %#v", remoteDescriptor)

	logging.Info("Validating remote descriptor...")
	err = remoteDescriptor.Validate()
	if err != nil {
		logging.Warning(err.Error())
		return nil
	}

	logging.Notice("Remote descriptor valid")
	return remoteDescriptor
}

func chooseReferenceDescriptor(remoteDescriptor *apps.AppDescriptor, localDescriptor *apps.AppDescriptor) (referenceDescriptor *apps.AppDescriptor, err error) {
	if remoteDescriptor == nil && localDescriptor == nil {
		return nil, fmt.Errorf("Cannot run the application: it is not installed and cannot be downloaded")
	}

	if remoteDescriptor == nil {
		if localDescriptor.SkipUpdateCheck {
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

	if remoteDescriptor.Version.NewerThan(localDescriptor.Version) {
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

func tryToSaveReferenceDescriptor(referenceDescriptorCopy apps.AppDescriptor, localDescriptorPath string, originalBaseURL *url.URL) (referenceDescriptorSaved bool) {
	referenceDescriptorCopy.BaseURL = originalBaseURL

	logging.Info("Saving the reference descriptor as the local descriptor...")
	referenceDescriptorBytes, err := referenceDescriptorCopy.ToBytes()
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
