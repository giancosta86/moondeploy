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
	"github.com/op/go-logging"

	"github.com/giancosta86/moondeploy/v3/apps"
	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/launchers"
	"github.com/giancosta86/moondeploy/v3/log"
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
func Run(
	launcher launchers.Launcher,
	userInterface ui.UserInterface,
	bootDescriptor descriptors.AppDescriptor) (err error) {

	settings := launcher.GetSettings()

	//----------------------------------------------------------------------------

	setupUserInterface(launcher, userInterface)
	defer dismissUserInterface(userInterface)

	//----------------------------------------------------------------------------

	userInterface.SetHeader("Performing startup operations")

	log.Info("The boot descriptor is: %#v", bootDescriptor)

	//----------------------------------------------------------------------------

	userInterface.SetApp(bootDescriptor.GetName())

	//----------------------------------------------------------------------------

	appGallery := apps.NewAppGallery(settings.GetGalleryDirectory())
	log.Info("The app gallery is: %#v", appGallery)

	//----------------------------------------------------------------------------

	log.Info("Resolving the app...")
	app, err := appGallery.GetApp(bootDescriptor)
	if err != nil {
		return err
	}
	log.Notice("The app directory is: '%v'", app.Directory)

	log.Info("App is: %#v", app)

	firstRun := !app.DirectoryExists()
	log.Notice("Is this a first run for the app? %v", firstRun)

	//----------------------------------------------------------------------------

	if firstRun {
		log.Info("Now asking the user if the app can run...")

		canRun := app.CanPerformFirstRun(userInterface)
		if !canRun {
			return &ExecutionCanceled{}
		}

		log.Notice("The user agreed to proceed")

		log.Info("Ensuring the app dir is available...")
		err = app.EnsureDirectory()
		if err != nil {
			return err
		}
		log.Notice("App dir available")
	}

	//----------------------------------------------------------------------------

	log.Info("Locking the app dir...")
	err = app.LockDirectory()
	if err != nil {
		return err
	}
	defer func() {
		unlockErr := app.UnlockDirectory()
		if unlockErr != nil {
			log.Warning(unlockErr.Error())
		}
	}()

	log.Notice("App dir locked")

	//----------------------------------------------------------------------------

	log.Info("Resolving the local descriptor...")
	localDescriptor := app.GetLocalDescriptor()

	startedWithLocalDescriptor := localDescriptor != nil
	log.Info("Started with local descriptor? %v", startedWithLocalDescriptor)

	if startedWithLocalDescriptor {
		log.Info("Checking that local descriptor and boot descriptor actually match...")
		err = descriptors.CheckDescriptorMatch(localDescriptor, bootDescriptor)
		if err != nil {
			return err
		}
		log.Notice("The descriptors match correctly")
	}

	//----------------------------------------------------------------------------

	log.Info("Resolving the remote descriptor...")
	remoteDescriptor := app.GetRemoteDescriptor()

	if remoteDescriptor != nil {
		log.Info("Checking that remote descriptor and boot descriptor actually match...")
		err = descriptors.CheckDescriptorMatch(remoteDescriptor, bootDescriptor)
		if err != nil {
			return err
		}
		log.Notice("The descriptors match correctly")
	}

	//----------------------------------------------------------------------------

	log.Info("Now choosing the reference descriptor...")
	referenceDescriptor, err := app.GetReferenceDescriptor()
	if err != nil {
		return err
	}

	log.Info("The reference descriptor is: %#v", referenceDescriptor)

	//----------------------------------------------------------------------------

	err = referenceDescriptor.CheckRequirements()
	if err != nil {
		return err
	}

	//----------------------------------------------------------------------------

	userInterface.SetApp(referenceDescriptor.GetTitle())

	//----------------------------------------------------------------------------

	log.Info("Resolving the OS-specific app command line...")
	commandLine := referenceDescriptor.GetCommandLine()
	log.Info("Command line resolved: %v", commandLine)

	//----------------------------------------------------------------------------

	err = app.CheckFiles(settings, userInterface)
	if err != nil {
		return err
	}

	//----------------------------------------------------------------------------

	userInterface.SetHeader("Preparing the command...")

	command := app.PrepareCommand(commandLine)
	log.Notice("Command created")

	log.Info("Command path: %v", command.Path)
	log.Info("Command arguments: %v", command.Args)

	//----------------------------------------------------------------------------

	referenceDescriptorSaved := app.SaveReferenceDescriptor()

	if !startedWithLocalDescriptor && referenceDescriptorSaved {
		if userInterface.AskForDesktopShortcut(referenceDescriptor) {
			log.Info("Creating desktop shortcut...")

			err = app.CreateDesktopShortcut(launcher, referenceDescriptor)
			if err != nil {
				log.Warning("Could not create desktop shortcut: %v", err)
			} else {
				log.Notice("Desktop shortcut created")
			}
		} else {
			log.Info("The user refused to create a desktop shortcut")
		}
	}

	//----------------------------------------------------------------------------

	app.UnlockDirectory()

	//----------------------------------------------------------------------------

	userInterface.SetHeader("Launching the application")
	userInterface.SetStatus("")

	return app.Launch(command, settings, userInterface)
}

func setupUserInterface(launcher launchers.Launcher, userInterface ui.UserInterface) {
	userInterface.SetApp(launcher.GetTitle())

	log.SetCallback(func(level logging.Level, message string) {
		if level <= logging.INFO {
			userInterface.SetStatus(message)
		}
	})

	userInterface.Show()
}

func dismissUserInterface(userInterface ui.UserInterface) {
	log.SetCallback(func(level logging.Level, message string) {})

	userInterface.Hide()
}
