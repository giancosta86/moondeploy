/*§
  ===========================================================================
  MoonDeploy
  ===========================================================================
  Copyright (C) 2015 Gianluca Costa
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

package main

import (
	"fmt"
	"os"

	"github.com/giancosta86/moondeploy"
	"github.com/giancosta86/moondeploy/custom"
	"github.com/giancosta86/moondeploy/engine"
	"github.com/giancosta86/moondeploy/logging"

	"github.com/giancosta86/moondeploy/client/verbs"
)

func main() {
	fmt.Println(moondeploy.Title)

	if len(os.Args) < 2 {
		exitWithUsage()
	}

	settings := loadSettings()

	setLoggingLevel(settings)

	command := os.Args[1]
	err := executeCommand(command, settings)

	switch err.(type) {
	case nil:
		os.Exit(moondeploy.ExitCodeSuccess)

	case *engine.ExecutionCanceled:
		exitWithCancel()

	case *verbs.InvalidCommandLineArguments:
		exitWithUsage()

	default:
		exitWithError(err)
	}
}

func loadSettings() (settings *custom.Settings) {
	logging.Info("Loading settings...")

	settings, err := custom.ComputeUserSettings()
	if err != nil {
		exitWithError(err)
	}

	logging.Notice("Settings loaded: %#v", settings)

	return settings
}

func setLoggingLevel(settings *custom.Settings) {
	logging.Info("Configuring the logging level...")
	loggingLevel := settings.GetLoggingLevel()
	logging.Notice("Requested logging level: %v", loggingLevel)
	logging.SetLevel(loggingLevel)
	logging.Notice("Logging level set")
}

func executeCommand(command string, settings *custom.Settings) (err error) {
	switch command {
	case moondeploy.ServeVerb:
		return verbs.DoServe()

	default:
		return verbs.DoRun(settings)
	}
}

func exitWithCancel() {
	logging.Warning("*** EXECUTION CANCELED ***")
	os.Exit(moondeploy.ExitCodeCanceled)
}

func exitWithError(err error) {
	logging.Error(err.Error())
	os.Exit(moondeploy.ExitCodeError)
}

func exitWithUsage() {
	fmt.Println()
	fmt.Println()
	fmt.Printf("USAGE: <%v> <app descriptor file>|(<command> <parameters>)\n", os.Args[0])
	fmt.Println()
	fmt.Println("Available commands")
	fmt.Println()
	fmt.Printf("%v <port> <directory>\n", moondeploy.ServeVerb)
	fmt.Println("\tStarts an HTTP server on <port> serving files from <directory>")
	fmt.Println()

	os.Exit(moondeploy.ExitCodeError)
}
