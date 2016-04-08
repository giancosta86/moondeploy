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

package termui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/giancosta86/caravel/terminals"

	"github.com/giancosta86/moondeploy/v3"
	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/launchers"
	"github.com/giancosta86/moondeploy/v3/log"
	"github.com/giancosta86/moondeploy/v3/ui"
)

type TerminalUserInterface struct {
	terminal terminals.Terminal

	app      string
	header   string
	status   string
	progress float64
}

func NewTerminalUserInterface(launcher launchers.Launcher, terminal terminals.Terminal) *TerminalUserInterface {
	log.Debug("Terminal rows: %v", terminal.GetRows())
	log.Debug("Terminal columns: %v", terminal.GetColumns())

	return &TerminalUserInterface{
		terminal: terminal,
	}
}

func (userInterface *TerminalUserInterface) ShowError(message string) {
	userInterface.Hide()

	log.Error(message)
}

func (userInterface *TerminalUserInterface) askYesNo(prompt string) (yesResponse bool) {
	terminal := userInterface.terminal

	for {
		terminal.ResetStyle()
		userInterface.setupColors()
		terminal.Clear()
		terminal.ShowCursor()

		userInterface.drawTitle()

		terminal.MoveCursor(8, 1)
		fmt.Printf("%v [Y/N]: ", prompt)

		reader := bufio.NewReader(os.Stdin)

		userInput, err := reader.ReadString('\n')
		if err != nil {
			userInterface.ShowError(err.Error())
			os.Exit(v3.ExitCodeError)
		}

		userInput = strings.TrimSpace(userInput)

		switch strings.ToUpper(userInput) {
		case "Y":
			return true
		case "N":
			return false
		}
	}
}

func (userInterface *TerminalUserInterface) styleFirstRunPrompt(prompt string) string {
	if !userInterface.terminal.SupportsANSI() {
		return prompt
	}

	prompt = strings.Replace(prompt, "\n\nTitle:", "\n\n\033[1m   Title:\033[21m", -1)
	prompt = strings.Replace(prompt, "\n\nPublisher:", "\n\n\033[1m   Publisher:\033[21m", -1)
	prompt = strings.Replace(prompt, "\n\nAddress:", "\n\n\033[1m   Address:\033[21m", -1)
	prompt = strings.Replace(prompt, "\n\nWARNING:", "\n\n\033[1mWARNING:\033[21m", -1)

	return prompt
}

func (userInterface *TerminalUserInterface) AskForSecureFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool) {
	prompt := userInterface.styleFirstRunPrompt(
		ui.FormatSecureFirstRunPrompt(bootDescriptor))

	return userInterface.askYesNo(prompt)
}

func (userInterface *TerminalUserInterface) AskForUntrustedFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool) {
	prompt := userInterface.styleFirstRunPrompt(
		ui.FormatUntrustedFirstRunPrompt(bootDescriptor))

	return userInterface.askYesNo(prompt)
}

func (userInterface *TerminalUserInterface) SetApp(app string) {
	userInterface.app = app
	userInterface.redraw()
}

func (userInterface *TerminalUserInterface) SetHeader(header string) {
	userInterface.header = header
	userInterface.redraw()
}

func (userInterface *TerminalUserInterface) SetStatus(status string) {
	userInterface.status = status
	userInterface.redraw()
}

func (userInterface *TerminalUserInterface) SetProgress(progress float64) {
	userInterface.progress = progress
	userInterface.redraw()
}

func (userInterface *TerminalUserInterface) AskForDesktopShortcut(referenceDescriptor descriptors.AppDescriptor) (canCreate bool) {
	prompt := ui.FormatDesktopShortcutPrompt(referenceDescriptor)
	return userInterface.askYesNo(prompt)
}

func (userInterface *TerminalUserInterface) Show() {

}

func (userInterface *TerminalUserInterface) Hide() {
	terminal := userInterface.terminal

	terminal.ResetStyle()
	terminal.ShowCursor()
	terminal.Clear()
}

func (userInterface *TerminalUserInterface) setupColors() {
	terminal := userInterface.terminal

	terminal.SetBackgroundColor(195)
	terminal.SetForegroundColor(22)
}

func (userInterface *TerminalUserInterface) redraw() {
	terminal := userInterface.terminal

	terminal.ResetStyle()
	userInterface.setupColors()
	terminal.Clear()

	userInterface.drawTitle()

	terminal.EnableTextBold()
	terminal.PrintCenteredInRow(8, userInterface.header)
	terminal.DisableTextBold()

	terminal.MoveCursor(12, 2)
	fmt.Print(userInterface.status)

	if 0 < userInterface.progress && userInterface.progress < 1 {
		terminal.DrawHorizontalProgressBar(16, 2, terminal.GetColumns()-20, userInterface.progress)
	}

	terminal.HideCursor()
	terminal.EnableTextHidden()

	//TODO: DEL THIS
	time.Sleep(500 * time.Millisecond)
}

func (userInterface *TerminalUserInterface) drawTitle() {
	if userInterface.app == "" {
		return
	}

	terminal := userInterface.terminal

	solidLine := strings.Repeat("-", terminal.GetColumns())

	terminal.MoveCursor(1, 1)
	fmt.Print(solidLine)

	terminal.EnableTextBold()
	terminal.PrintCenteredInRow(3, strings.ToUpper(userInterface.app))
	terminal.DisableTextBold()

	terminal.MoveCursor(5, 1)
	fmt.Print(solidLine)
}
