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

	"github.com/giancosta86/moondeploy/v3"
	"github.com/giancosta86/moondeploy/v3/engine"
	"github.com/giancosta86/moondeploy/v3/launchers"
	"github.com/giancosta86/moondeploy/v3/log"

	"github.com/giancosta86/moondeploy/moonclient/verbs"
)

func main() {
	launcher := getMoonLauncher()
	fmt.Println(launcher.GetTitle())

	if len(os.Args) < 2 {
		exitWithUsage()
	}

	settings := launcher.GetSettings()

	initializeLogging(settings)

	log.Info("Launcher is: %#v", launcher)

	command := os.Args[1]
	err := executeCommand(launcher, command)

	switch err.(type) {
	case nil:
		os.Exit(v3.ExitCodeSuccess)

	case *engine.ExecutionCanceled:
		exitWithCancel()

	case *verbs.InvalidCommandLineArguments:
		exitWithUsage()

	default:
		exitWithError(err)
	}
}

func executeCommand(launcher launchers.Launcher, command string) (err error) {
	settings := launcher.GetSettings()

	switch command {
	case verbs.Serve:
		return verbs.DoServe()

	default:
		return verbs.DoRun(launcher, settings)
	}
}

func exitWithCancel() {
	log.Warning("*** EXECUTION CANCELED ***")
	os.Exit(v3.ExitCodeCanceled)
}

func exitWithError(err error) {
	log.Error(err.Error())
	os.Exit(v3.ExitCodeError)
}

func exitWithUsage() {
	fmt.Println()
	fmt.Println()
	fmt.Printf("USAGE: <%v> <app descriptor file>|(<command> <parameters>)\n", os.Args[0])
	fmt.Println()
	fmt.Println("Available commands")
	fmt.Println()
	fmt.Printf("%v <port> <directory>\n", verbs.Serve)
	fmt.Println("\tStarts an HTTP server on <port> serving files from <directory>")
	fmt.Println()

	os.Exit(v3.ExitCodeError)
}
