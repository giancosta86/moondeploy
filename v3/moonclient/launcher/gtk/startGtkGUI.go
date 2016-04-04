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
	"github.com/giancosta86/moondeploy/v3/custom"
	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/engine"
	"github.com/giancosta86/moondeploy/v3/log"
	"github.com/giancosta86/moondeploy/v3/moonclient"
	"github.com/giancosta86/moondeploy/v3/ui/gtkui"
	"github.com/gotk3/gotk3/gtk"
	"github.com/op/go-logging"
)

type asyncResult struct {
	userInterface *gtkui.GtkUserInterface
	err           error
}

func StartGUI(bootDescriptorPath string, settings *custom.Settings) (err error) {
	log.Info("Initializing GTK...")
	gtkui.InitGTK()
	log.Notice("GTK initialized")

	resultChannel := make(chan asyncResult)
	defer close(resultChannel)

	go backgroundCollector(bootDescriptorPath, settings, resultChannel)

	log.Info("Starting GTK main loop...")
	gtk.Main()

	log.SetCallback(func(level logging.Level, message string) {})
	log.Notice("GTK main loop terminated")

	select {
	case result := <-resultChannel:
		err = result.err
		log.Info("Result retrieved from channel. Err is: '%v'", err)

		if result.userInterface != nil && result.userInterface.IsClosedByUser() {
			return &engine.ExecutionCanceled{}
		}

		if err != nil {
			return err
		}

		log.Notice("OK")
		return nil
	default:
		log.Info("The user has manually closed the program")
		return &engine.ExecutionCanceled{}
	}
}

func backgroundCollector(bootDescriptorPath string, settings *custom.Settings, resultChannel chan asyncResult) {
	result := backgroundProcessing(bootDescriptorPath, settings)
	userInterface := result.userInterface
	err := result.err

	log.SetCallback(func(level logging.Level, message string) {})
	log.Info("Result returned by the background routine. Is UI available? %v. Err is: '%v'", userInterface != nil, err)

	if err != nil && userInterface != nil {
		switch err.(type) {

		case *engine.ExecutionCanceled:
			break

		default:
			userInterface.ShowError(err.Error())
		}
	}

	log.Info("Now programmatically quitting GTK")
	gtk.MainQuit()

	resultChannel <- result
}

func backgroundProcessing(bootDescriptorPath string, settings *custom.Settings) asyncResult {
	log.Info("Creating the user interface...")
	userInterface, err := gtkui.NewGtkUserInterface()
	if err != nil {
		return asyncResult{
			userInterface: nil,
			err:           err,
		}
	}
	log.Notice("User interface created")

	showUserInterface(userInterface)

	//----------------------------------------------------------------------------
	log.Info("Opening boot descriptor: %v", bootDescriptorPath)

	bootDescriptor, err := descriptors.NewAppDescriptorFromPath(bootDescriptorPath)
	if err != nil {
		return asyncResult{
			userInterface: userInterface,
			err:           err,
		}
	}

	log.Notice("Boot descriptor ready")
	//----------------------------------------------------------------------------

	log.Info("Starting the launch process...")

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

func showUserInterface(userInterface *gtkui.GtkUserInterface) {
	userInterface.SetApp(moonclient.Title)
	userInterface.SetHeader("Loading the boot descriptor")

	log.Info("Registering user interface for log...")
	log.SetCallback(func(level logging.Level, message string) {
		if level <= logging.NOTICE {
			userInterface.SetStatus(message)
		}
	})
	log.Notice("User interface registered")

	log.Info("Showing the loading dialog...")
	userInterface.ShowLoader()
	log.Notice("Loading dialog shown")
}
