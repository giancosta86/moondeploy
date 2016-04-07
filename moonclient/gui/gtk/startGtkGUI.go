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

package gtk

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/op/go-logging"

	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/engine"
	"github.com/giancosta86/moondeploy/v3/launchers"
	"github.com/giancosta86/moondeploy/v3/log"
	"github.com/giancosta86/moondeploy/v3/ui/gtkui"
)

type guiOutcomeStruct struct {
	userInterface *gtkui.GtkUserInterface
	err           error
}

func StartGUI(launcher launchers.Launcher, bootDescriptorPath string) (err error) {
	log.Info("Initializing GTK...")
	gtkui.InitGTK()
	log.Notice("GTK initialized")

	guiOutcomeChannel := make(chan guiOutcomeStruct)
	defer close(guiOutcomeChannel)

	go backgroundOrchestrator(launcher, bootDescriptorPath, guiOutcomeChannel)

	log.Info("Starting GTK main loop...")
	gtk.Main()

	log.SetCallback(func(level logging.Level, message string) {})
	log.Notice("GTK main loop terminated")

	select {
	case guiOutcome := <-guiOutcomeChannel:
		log.Info("Outcome retrieved from the GUI channel")

		if guiOutcome.userInterface != nil && guiOutcome.userInterface.IsClosedByUser() {
			return &engine.ExecutionCanceled{}
		}

		err = guiOutcome.err
		if err != nil {
			log.Warning("Err is: %v", err)
			return err
		}

		log.Notice("OK")
		return nil
	default:
		log.Info("The user has manually closed the program")
		return &engine.ExecutionCanceled{}
	}
}

func backgroundOrchestrator(launcher launchers.Launcher, bootDescriptorPath string, guiOutcomeChannel chan guiOutcomeStruct) {
	outcome := runEngineWithGtk(launcher, bootDescriptorPath)
	userInterface := outcome.userInterface
	err := outcome.err

	log.SetCallback(func(level logging.Level, message string) {})
	log.Info("Result returned by the background routine. Is UI available? %v", userInterface != nil)

	if err != nil {
		log.Warning("Err is: %v", err)

		if userInterface != nil {
			switch err.(type) {

			case *engine.ExecutionCanceled:
				break

			default:
				userInterface.ShowError(err.Error())
			}
		}
	}

	log.Info("Now programmatically quitting GTK")
	gtk.MainQuit()

	guiOutcomeChannel <- outcome
}

func runEngineWithGtk(launcher launchers.Launcher, bootDescriptorPath string) guiOutcomeStruct {
	log.Info("Creating the GTK+ user interface...")

	userInterface, err := gtkui.NewGtkUserInterface(launcher)
	if err != nil {
		return guiOutcomeStruct{
			userInterface: nil,
			err:           err,
		}
	}

	log.Notice("User interface created")

	//----------------------------------------------------------------------------
	log.Info("Opening boot descriptor: %v", bootDescriptorPath)

	bootDescriptor, err := descriptors.NewAppDescriptorFromPath(bootDescriptorPath)
	if err != nil {
		return guiOutcomeStruct{
			userInterface: userInterface,
			err:           err,
		}
	}

	log.Notice("Boot descriptor ready")
	//----------------------------------------------------------------------------

	log.Info("Starting the launch process...")

	err = engine.Run(launcher, userInterface, bootDescriptor)
	return guiOutcomeStruct{
		userInterface: userInterface,
		err:           err,
	}
}
