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

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/giancosta86/moondeploy/v3/config"
	"github.com/giancosta86/moondeploy/v3/engine"
	"github.com/giancosta86/moondeploy/v3/log"
	"github.com/giancosta86/moondeploy/v3/moonclient"

	"github.com/giancosta86/moondeploy/v3/moonclient/verbs"
)

const ExitCodeSuccess = 0
const ExitCodeError = 1
const ExitCodeCanceled = 2

const ServeVerb = "serve"

func main() {
	fmt.Println(moonclient.Title)

	if len(os.Args) < 2 {
		exitWithUsage()
	}

	logsDirectory := filepath.Join(moonclient.Directory, "logs")
	log.SwitchToFile(logsDirectory)

	settings := loadSettings()

	setLoggingLevel(settings)

	command := os.Args[1]
	err := executeCommand(command, settings)

	switch err.(type) {
	case nil:
		os.Exit(ExitCodeSuccess)

	case *engine.ExecutionCanceled:
		exitWithCancel()

	case *verbs.InvalidCommandLineArguments:
		exitWithUsage()

	default:
		exitWithError(err)
	}
}

func loadSettings() (settings *config.Settings) {
	log.Info("Loading settings...")

	settings, err := moonclient.ComputeUserSettings()
	if err != nil {
		exitWithError(err)
	}

	log.Info("Settings loaded: %#v", settings)

	return settings
}

func setLoggingLevel(settings *config.Settings) {
	log.Info("Configuring the logging level...")
	loggingLevel := settings.GetLoggingLevel()
	log.Notice("Requested logging level: %v", loggingLevel)
	log.SetLevel(loggingLevel)
	log.Notice("Logging level set")
}

func executeCommand(command string, settings *config.Settings) (err error) {
	switch command {
	case ServeVerb:
		return verbs.DoServe()

	default:
		return verbs.DoRun(settings)
	}
}

func exitWithCancel() {
	log.Warning("*** EXECUTION CANCELED ***")
	os.Exit(ExitCodeCanceled)
}

func exitWithError(err error) {
	log.Error(err.Error())
	os.Exit(ExitCodeError)
}

func exitWithUsage() {
	fmt.Println()
	fmt.Println()
	fmt.Printf("USAGE: <%v> <app descriptor file>|(<command> <parameters>)\n", os.Args[0])
	fmt.Println()
	fmt.Println("Available commands")
	fmt.Println()
	fmt.Printf("%v <port> <directory>\n", ServeVerb)
	fmt.Println("\tStarts an HTTP server on <port> serving files from <directory>")
	fmt.Println()

	os.Exit(ExitCodeError)
}
