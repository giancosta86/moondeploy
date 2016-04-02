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

package verbs

import (
	"os"

	"github.com/gotk3/gotk3/gtk"

	"github.com/giancosta86/moondeploy/v3/apps"
	"github.com/giancosta86/moondeploy/v3/custom"
	"github.com/giancosta86/moondeploy/v3/engine"
	"github.com/giancosta86/moondeploy/v3/logging"
	"github.com/giancosta86/moondeploy/v3/moonclient"
	"github.com/giancosta86/moondeploy/v3/ui/gtkui"
)

type asyncResult struct {
	userInterface *gtkui.GtkUserInterface
	err           error
}

func DoRun(settings *custom.Settings) (err error) {
	bootDescriptorPath := os.Args[1]

	logging.Info("Initializing GTK...")
	gtkui.InitGTK()
	logging.Notice("GTK initialized")

	resultChannel := make(chan asyncResult)
	defer close(resultChannel)

	go backgroundCollector(bootDescriptorPath, settings, resultChannel)

	logging.Info("Starting GTK main loop...")
	gtk.Main()

	logging.SetCallback(func(message string) {})
	logging.Notice("GTK main loop terminated")

	select {
	case result := <-resultChannel:
		err = result.err
		logging.Info("Result retrieved from channel. Err is: '%v'", err)

		if result.userInterface != nil && result.userInterface.IsClosedByUser() {
			return &engine.ExecutionCanceled{}
		}

		if err != nil {
			return err
		}

		logging.Notice("OK")
		return nil
	default:
		logging.Info("The user has manually closed the program")
		return &engine.ExecutionCanceled{}
	}
}

func backgroundCollector(bootDescriptorPath string, settings *custom.Settings, resultChannel chan asyncResult) {
	result := backgroundProcessing(bootDescriptorPath, settings)
	userInterface := result.userInterface
	err := result.err

	logging.SetCallback(func(message string) {})
	logging.Info("Result returned by the background routine. Is UI available? %v. Err is: '%v'", userInterface != nil, err)

	if err != nil && userInterface != nil {
		switch err.(type) {

		case *engine.ExecutionCanceled:
			break

		default:
			userInterface.ShowError(err.Error())
		}
	}

	logging.Info("Now programmatically quitting GTK")
	gtk.MainQuit()

	resultChannel <- result
}

func backgroundProcessing(bootDescriptorPath string, settings *custom.Settings) asyncResult {
	logging.Info("Creating the user interface...")
	userInterface, err := gtkui.NewGtkUserInterface()
	if err != nil {
		return asyncResult{
			userInterface: nil,
			err:           err,
		}
	}
	logging.Notice("User interface created")

	startUserInterface(userInterface)

	//----------------------------------------------------------------------------
	logging.Info("Opening boot descriptor: %v", bootDescriptorPath)

	bootDescriptor, err := apps.NewAppDescriptorFromPath(bootDescriptorPath)
	if err != nil {
		return asyncResult{
			userInterface: userInterface,
			err:           err,
		}
	}

	logging.Notice("Boot descriptor ready")
	//----------------------------------------------------------------------------

	logging.Info("Starting the launch process...")

	err = engine.Run(bootDescriptor, settings, userInterface)
	if err != nil {
		return asyncResult{
			userInterface: userInterface,
			err:           err,
		}
	}

	return asyncResult{
		userInterface: userInterface,
		err:           nil,
	}
}

func startUserInterface(userInterface *gtkui.GtkUserInterface) {
	userInterface.SetApp(moonclient.Title)
	userInterface.SetHeader("Loading the boot descriptor")

	logging.Info("Registering user interface for logging...")
	logging.SetCallback(func(message string) {
		userInterface.SetStatus(message)
	})
	logging.Notice("User interface registered")

	logging.Info("Showing the loading dialog...")
	userInterface.ShowLoader()
	logging.Notice("Loading dialog shown")
}
