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
	"github.com/giancosta86/moondeploy/v3/apps"
	"github.com/giancosta86/moondeploy/v3/custom"
	"github.com/giancosta86/moondeploy/v3/descriptors"
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
func Run(bootDescriptor descriptors.AppDescriptor, settings *custom.Settings, userInterface ui.UserInterface) (err error) {
	userInterface.SetHeader("Performing startup operations")

	logging.Info("The boot descriptor is: %#v", bootDescriptor)

	//----------------------------------------------------------------------------

	userInterface.SetApp(bootDescriptor.GetName())

	//----------------------------------------------------------------------------

	appGallery := apps.NewAppGallery(settings.GalleryDir)
	logging.Notice("The app gallery is: %#v", appGallery)

	//----------------------------------------------------------------------------

	logging.Info("Resolving the app...")
	app, err := appGallery.GetApp(bootDescriptor)
	if err != nil {
		return err
	}
	logging.Notice("App is: %#v", app)

	firstRun := !app.DirectoryExists()
	logging.Notice("Is this a first run for the app? %v", firstRun)

	//----------------------------------------------------------------------------

	if firstRun {
		logging.Info("Now asking the user if the app can run...")

		canRun := app.CanPerformFirstRun(userInterface)
		if !canRun {
			return &ExecutionCanceled{}
		}

		logging.Notice("The user agreed to proceed")

		logging.Info("Ensuring the app dir is available...")
		err = app.EnsureDirectory()
		if err != nil {
			return err
		}
		logging.Notice("App dir available")
	}

	//----------------------------------------------------------------------------

	logging.Info("Locking the app dir...")
	err = app.LockDirectory()
	if err != nil {
		return err
	}
	defer func() {
		unlockErr := app.UnlockDirectory()
		if unlockErr != nil {
			logging.Warning(unlockErr.Error())
		}
	}()

	logging.Notice("App dir locked")

	//----------------------------------------------------------------------------

	logging.Info("Resolving the local descriptor...")
	localDescriptor := app.GetLocalDescriptor()

	startedWithLocalDescriptor := localDescriptor != nil
	logging.Info("Started with local descriptor? %v", startedWithLocalDescriptor)

	if startedWithLocalDescriptor {
		logging.Info("Checking that local descriptor and boot descriptor actually match...")
		err = descriptors.CheckDescriptorMatch(localDescriptor, bootDescriptor)
		if err != nil {
			return err
		}
		logging.Notice("The descriptors match correctly")
	}

	//----------------------------------------------------------------------------

	logging.Info("Resolving the remote descriptor...")
	remoteDescriptor := app.GetRemoteDescriptor()

	if remoteDescriptor != nil {
		logging.Info("Checking that remote descriptor and boot descriptor actually match...")
		err = descriptors.CheckDescriptorMatch(remoteDescriptor, bootDescriptor)
		if err != nil {
			return err
		}
		logging.Notice("The descriptors match correctly")
	}

	//----------------------------------------------------------------------------

	logging.Info("Now choosing the reference descriptor...")
	referenceDescriptor, err := app.GetReferenceDescriptor()
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
	logging.Notice("Command line resolved: %v", commandLine)

	//----------------------------------------------------------------------------

	err = app.CheckFiles(settings, userInterface)
	if err != nil {
		return err
	}

	//----------------------------------------------------------------------------

	userInterface.SetHeader("Preparing the command...")

	command := app.PrepareCommand(commandLine)
	logging.Notice("Command created")

	logging.Info("Command path: %v", command.Path)
	logging.Info("Command arguments: %v", command.Args)

	//----------------------------------------------------------------------------

	referenceDescriptorSaved := app.SaveReferenceDescriptor()

	if !startedWithLocalDescriptor && referenceDescriptorSaved {
		if userInterface.AskForDesktopShortcut(referenceDescriptor) {
			logging.Info("Creating desktop shortcut...")

			err = app.CreateDesktopShortcut(referenceDescriptor)
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

	app.UnlockDirectory()

	//----------------------------------------------------------------------------

	userInterface.SetHeader("Launching the application")

	return app.Launch(command, settings, userInterface)
}
