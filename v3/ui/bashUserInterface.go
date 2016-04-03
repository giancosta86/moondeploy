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

package ui

import (
	"fmt"

	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/logging"
)

type BashUserInterface struct {
	app      string
	header   string
	status   string
	progress float64
}

func NewBashUserInterface() *BashUserInterface {
	return &BashUserInterface{}
}

func (userInterface *BashUserInterface) ShowError(message string) {
	logging.Error(message)
}

func (userInterface *BashUserInterface) AskForSecureFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool) {
	return true
}

func (userInterface *BashUserInterface) AskForUntrustedFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool) {
	return true
}

func (userInterface *BashUserInterface) SetApp(app string) {
	userInterface.app = app
	userInterface.redraw()
}

func (userInterface *BashUserInterface) SetHeader(header string) {
	userInterface.header = header
	userInterface.redraw()
}

func (userInterface *BashUserInterface) SetStatus(status string) {
	userInterface.status = status
	userInterface.redraw()
}

func (userInterface *BashUserInterface) SetProgress(progress float64) {
	userInterface.progress = progress
	userInterface.redraw()
}

func (userInterface *BashUserInterface) AskForDesktopShortcut(referenceDescriptor descriptors.AppDescriptor) (canCreate bool) {
	return true
}

func (userInterface *BashUserInterface) HideLoader() {
	resetTerminal()
	clearScreen()
}

func clearScreen() {
	fmt.Print("\033[2J")
	fmt.Print("\033[0;0H")
}

func resetTerminal() {
	fmt.Print("\033[0m")
}

func setupColors() {
	//Background
	fmt.Print("\033[48;5;17m")

	//Foreground
	fmt.Print("\033[38;5;186m")
}

func (userInterface *BashUserInterface) redraw() {
	clearScreen()

	setupColors()

	fmt.Println("================================================================================")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Println("|                                                                              |")
	fmt.Print("================================================================================")
}
