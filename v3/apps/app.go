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
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3/config"
	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/launchers"
	"github.com/giancosta86/moondeploy/v3/log"
	"github.com/giancosta86/moondeploy/v3/ui"
)

const filesDirName = "files"
const lockFileName = "App.lock"

type App struct {
	Directory string

	bootDescriptor descriptors.AppDescriptor

	filesDirectory string

	lockFile *os.File

	localDescriptor       descriptors.AppDescriptor
	localDescriptorCached bool
	localDescriptorPath   string

	remoteDescriptor       descriptors.AppDescriptor
	remoteDescriptorCached bool

	referenceDescriptor       descriptors.AppDescriptor
	referenceDescriptorCached bool
}

func (app *App) DirectoryExists() bool {
	return caravel.DirectoryExists(app.Directory)
}

func (app *App) CanPerformFirstRun(userInterface ui.UserInterface) bool {
	if caravel.IsSecureURL(app.bootDescriptor.GetDeclaredBaseURL()) {
		return userInterface.AskForSecureFirstRun(app.bootDescriptor)
	} else {
		return userInterface.AskForUntrustedFirstRun(app.bootDescriptor)
	}
}

func (app *App) EnsureDirectory() (err error) {
	err = os.MkdirAll(app.Directory, 0700)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) GetLocalDescriptor() (localDescriptor descriptors.AppDescriptor) {
	if app.localDescriptorCached {
		return app.localDescriptor
	}

	app.localDescriptorCached = true

	app.localDescriptorPath = filepath.Join(app.Directory, app.bootDescriptor.GetDescriptorFileName())

	if !caravel.FileExists(app.localDescriptorPath) {
		log.Notice("The local descriptor is missing")
		return nil
	}

	log.Notice("The local descriptor has been found! Deserializing...")
	localDescriptor, err := descriptors.NewAppDescriptorFromPath(app.localDescriptorPath)
	if err != nil {
		log.Warning(err.Error())
		return nil
	}
	log.Notice("Local descriptor deserialized")

	log.Info("The local descriptor is: %#v", localDescriptor)

	app.localDescriptor = localDescriptor
	app.localDescriptorPath = app.localDescriptorPath

	return localDescriptor
}

func (app *App) GetRemoteDescriptor() (remoteDescriptor descriptors.AppDescriptor) {
	if app.remoteDescriptorCached {
		return app.remoteDescriptor
	}

	app.remoteDescriptorCached = true

	bootDescriptor := app.bootDescriptor
	localDescriptor := app.GetLocalDescriptor()

	var remoteDescriptorURL *url.URL
	var err error

	if localDescriptor != nil {
		if localDescriptor.IsSkipUpdateCheck() {
			log.Notice("Skipping update check, as requested by the local descriptor")
			return nil
		} else {
			remoteDescriptorURL, err = localDescriptor.GetFileURL(localDescriptor.GetDescriptorFileName())
		}
	} else {
		remoteDescriptorURL, err = bootDescriptor.GetFileURL(bootDescriptor.GetDescriptorFileName())
	}

	if err != nil {
		log.Warning(err.Error())
		return nil
	}

	log.Notice("The remote descriptor's URL is: %v", remoteDescriptorURL)

	log.Info("Retrieving the remote descriptor...")
	remoteDescriptorBytes, err := caravel.RetrieveFromURL(remoteDescriptorURL)
	if err != nil {
		log.Warning(err.Error())
		return nil
	}
	log.Notice("Remote descriptor retrieved")

	log.Info("Deserializing the remote descriptor...")
	remoteDescriptor, err = descriptors.NewAppDescriptorFromBytes(remoteDescriptorBytes)
	if err != nil {
		log.Warning(err.Error())
		return nil
	}
	log.Notice("Remote descriptor deserialized")

	log.Info("The remote descriptor is: %#v", remoteDescriptor)

	app.remoteDescriptor = remoteDescriptor

	return remoteDescriptor
}

func (app *App) GetReferenceDescriptor() (referenceDescriptor descriptors.AppDescriptor, err error) {
	if app.referenceDescriptorCached {
		return app.referenceDescriptor, nil
	}

	app.referenceDescriptorCached = true

	localDescriptor := app.GetLocalDescriptor()
	remoteDescriptor := app.GetRemoteDescriptor()

	if remoteDescriptor == nil && localDescriptor == nil {
		return nil, fmt.Errorf("Cannot run the application: it is not installed and cannot be downloaded")
	}

	if remoteDescriptor == nil {
		if localDescriptor.IsSkipUpdateCheck() {
			log.Info("The remote descriptor is missing as requested, so the local descriptor will be used")
		} else {
			log.Warning("The remote descriptor is missing, so the local descriptor will be used")
		}
		app.referenceDescriptor = localDescriptor
	} else if localDescriptor == nil {
		log.Notice("The local descriptor is missing, so the remote descriptor will be used")
		app.referenceDescriptor = remoteDescriptor
	} else if remoteDescriptor.GetAppVersion().NewerThan(localDescriptor.GetAppVersion()) {
		log.Notice("Switching to the remote descriptor, as it is more recent")
		app.referenceDescriptor = remoteDescriptor
	} else {
		log.Notice("Keeping the local descriptor, as the remote descriptor is NOT more recent")
		app.referenceDescriptor = localDescriptor
	}

	return app.referenceDescriptor, nil
}

func (app *App) PrepareCommand(commandLine []string) (command *exec.Cmd) {
	if caravel.DirectoryExists(app.filesDirectory) {
		os.Chdir(app.filesDirectory)
		log.Notice("Files directory set as the current directory")
	} else {
		os.Chdir(app.Directory)
		log.Notice("App directory set as the current directory")
	}

	log.Info("Creating the command...")

	if len(commandLine) == 1 {
		return exec.Command(commandLine[0])
	}

	return exec.Command(commandLine[0], commandLine[1:]...)
}

func (app *App) SaveReferenceDescriptor() (referenceDescriptorSaved bool) {
	referenceDescriptor, err := app.GetReferenceDescriptor()
	if err != nil {
		log.Warning("Cannot save the reference descriptor: %v", err)
		return false
	}

	log.Info("Saving the reference descriptor as the local descriptor...")
	referenceDescriptorBytes, err := referenceDescriptor.GetBytes()
	if err != nil {
		log.Error("Could not serialize the reference descriptor: %v", err)
		return false
	}

	err = ioutil.WriteFile(app.localDescriptorPath, referenceDescriptorBytes, 0600)
	if err != nil {
		log.Error("Could not save the reference descriptor: %v", err)
		return false
	}

	log.Notice("Reference descriptor saved")
	return true
}

func (app *App) Launch(command *exec.Cmd, settings config.Settings, userInterface ui.UserInterface) (err error) {
	log.Info("Starting the app...")

	log.Info("Hiding the user interface...")
	userInterface.HideLoader()
	log.Notice("User interface hidden")

	if settings.IsSkipAppOutput() {
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

func (app *App) GetActualIconPath(launcher launchers.Launcher) string {
	referenceDescriptor, err := app.GetReferenceDescriptor()
	if err != nil {
		log.Warning("Error while retrieving the reference descriptor: %v", err)
		return launcher.GetIconPath()
	}

	providedIconPath := referenceDescriptor.GetIconPath()

	if providedIconPath != "" {
		return filepath.Join(app.filesDirectory, providedIconPath)
	}

	return launcher.GetIconPath()
}
