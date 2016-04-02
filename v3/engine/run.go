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
	"path/filepath"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3/apps"
	"github.com/giancosta86/moondeploy/v3/custom"
	"github.com/giancosta86/moondeploy/v3/logging"
	"github.com/giancosta86/moondeploy/v3/ui"
)

/*
ExecutionCanceled is returned when the user explicitly interrupts the execution process,
for example refusing to install the application or closing the loading dialog.
*/
type ExecutionCanceled struct{}

func (err *ExecutionCanceled) Error() string {
	return "Execution canceled"
}

/*
Run is the entry point you must employ to create a custom installer, for example to
employ custom settings or a brand-new user interface, based on any technology
*/
func Run(bootDescriptor apps.AppDescriptor, settings *custom.Settings, userInterface ui.UserInterface) (err error) {
	userInterface.SetHeader("Performing startup operations")

	logging.Info("The boot descriptor is: %#v", bootDescriptor)

	//----------------------------------------------------------------------------

	userInterface.SetApp(bootDescriptor.GetName())

	//----------------------------------------------------------------------------

	appGalleryDir := settings.GalleryDir
	logging.Notice("The app gallery dir is: %v", appGalleryDir)

	//----------------------------------------------------------------------------

	logging.Info("Resolving the app dir...")
	appDir, err := resolveAppDir(bootDescriptor, appGalleryDir)
	if err != nil {
		return err
	}
	logging.Notice("App dir is: %v", appDir)

	appFilesDir := filepath.Join(appDir, apps.FilesDirName)
	logging.Info("App files dir is: %v", appFilesDir)

	firstRun := !caravel.DirectoryExists(appDir)
	logging.Notice("Is this a first run for the app? %v", firstRun)

	//----------------------------------------------------------------------------

	if firstRun {
		logging.Info("Now asking the user if the app can run...")
		err = ensureFirstRun(bootDescriptor, appDir, userInterface)
		if err != nil {
			return err
		}
	}

	//----------------------------------------------------------------------------

	logging.Info("Locking the app dir...")
	lockFile, err := lockAppDir(appDir)
	if err != nil {
		return err
	}
	defer func() {
		unlockErr := unlockAppDir(lockFile)
		if unlockErr != nil {
			logging.Warning(unlockErr.Error())
		}
	}()

	logging.Notice("App dir locked")

	//----------------------------------------------------------------------------

	logging.Info("Resolving the local descriptor...")
	localDescriptorPath := filepath.Join(appDir, bootDescriptor.GetDescriptorFileName())
	localDescriptor := getLocalDescriptor(localDescriptorPath)

	startedWithLocalDescriptor := localDescriptor != nil
	logging.Info("Started with local descriptor? %v", startedWithLocalDescriptor)

	if startedWithLocalDescriptor {
		logging.Info("Checking that local descriptor and boot descriptor actually match...")
		err = localDescriptor.CheckMatch(bootDescriptor)
		if err != nil {
			return err
		}
		logging.Notice("The descriptors match correctly")
	}

	//----------------------------------------------------------------------------

	var remoteDescriptor apps.AppDescriptor
	if localDescriptor != nil && localDescriptor.IsSkipUpdateCheck() {
		logging.Notice("Skipping update check, as requested by the local descriptor")
		remoteDescriptor = nil
	} else {
		logging.Info("Resolving the remote descriptor...")
		remoteDescriptor = getRemoteDescriptor(bootDescriptor, localDescriptor, userInterface)

		if remoteDescriptor != nil {
			logging.Info("Checking that remote descriptor and boot descriptor actually match...")
			err = remoteDescriptor.CheckMatch(bootDescriptor)
			if err != nil {
				return err
			}
			logging.Notice("The descriptors match correctly")
		}
	}

	//----------------------------------------------------------------------------

	logging.Info("Now choosing the reference descriptor...")
	referenceDescriptor, err := chooseReferenceDescriptor(remoteDescriptor, localDescriptor)
	if err != nil {
		return err
	}

	logging.Info("The reference descriptor is: %#v", referenceDescriptor)

	//----------------------------------------------------------------------------

	err = referenceDescriptor.CheckRequirements()
	if err != nil {
		return err
	}

	//----------------------------------------------------------------------------

	userInterface.SetApp(referenceDescriptor.GetTitle())

	//----------------------------------------------------------------------------

	logging.Info("Resolving the OS-specific app command line...")
	commandLine := referenceDescriptor.GetCommandLine()

	if len(commandLine) == 0 {
		return fmt.Errorf("Empty command line found")
	}
	logging.Notice("Command line resolved: %v", commandLine)

	//----------------------------------------------------------------------------

	if remoteDescriptor != nil {
		userInterface.SetHeader("Checking the app files")

		err = checkAppFiles(remoteDescriptor, localDescriptorPath, localDescriptor, appFilesDir, settings, userInterface)
		if err != nil {
			return err
		}
		logging.Notice("App files checked")
	}

	//----------------------------------------------------------------------------

	userInterface.SetHeader("Preparing the command...")

	command := prepareCommand(appDir, appFilesDir, commandLine)
	logging.Notice("Command created")

	logging.Info("Command path: %v", command.Path)
	logging.Info("Command arguments: %v", command.Args)

	//----------------------------------------------------------------------------

	referenceDescriptorSaved := tryToSaveReferenceDescriptor(referenceDescriptor, localDescriptorPath)

	if !startedWithLocalDescriptor && referenceDescriptorSaved {
		if userInterface.AskForDesktopShortcut(referenceDescriptor) {
			logging.Info("Creating desktop shortcut...")

			err = createDesktopShortcut(appFilesDir, localDescriptorPath, referenceDescriptor)
			if err != nil {
				logging.Warning("Could not create desktop shortcut: %v", err)
			} else {
				logging.Notice("Desktop shortcut created")
			}
		} else {
			logging.Info("The user refused to create a desktop shortcut")
		}
	}

	//----------------------------------------------------------------------------

	unlockAppDir(lockFile)

	//----------------------------------------------------------------------------

	userInterface.SetHeader("Launching the application")

	return launchApp(command, settings, userInterface)
}
