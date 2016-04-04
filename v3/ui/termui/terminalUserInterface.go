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

	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/log"
	"github.com/giancosta86/moondeploy/v3/ui"

	"github.com/op/go-logging"
)

type TerminalUserInterface struct {
	app      string
	header   string
	status   string
	progress float64

	terminal terminals.Terminal
}

func NewTerminalUserInterface(terminal terminals.Terminal) *TerminalUserInterface {
	log.Notice("Terminal rows: %v", terminal.GetRows())
	log.Notice("Terminal columns: %v", terminal.GetColumns())

	return &TerminalUserInterface{
		terminal: terminal,
	}
}

func (userInterface *TerminalUserInterface) ShowError(message string) {
	userInterface.HideLoader()

	log.Error(message)
}

func (userInterface *TerminalUserInterface) askYesNo(prompt string) (yesResponse bool) {
	terminal := userInterface.terminal
	for {
		terminal.ResetStyle()
		userInterface.setupColors()
		terminal.Clear()

		userInterface.drawTitle()

		terminal.ShowCursor()

		terminal.MoveCursor(8, 1)
		fmt.Printf("%v [Y/N]: ", prompt)

		reader := bufio.NewReader(os.Stdin)
		userInput, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
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

	prompt = strings.Replace(prompt, "\nTitle:", "\n\033[1m   Title:\033[21m", -1)
	prompt = strings.Replace(prompt, "\nPublisher:", "\n\033[1m   Publisher:\033[21m", -1)
	prompt = strings.Replace(prompt, "\nAddress:", "\n\033[1m   Address:\033[21m", -1)
	prompt = strings.Replace(prompt, "\nWARNING:", "\n\033[1mWARNING:\033[21m", -1)

	return prompt
}

func (userInterface *TerminalUserInterface) AskForSecureFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool) {
	prompt := ui.FormatSecureFirstRunPrompt(bootDescriptor)

	prompt = userInterface.styleFirstRunPrompt(prompt)

	return userInterface.askYesNo(prompt)
}

func (userInterface *TerminalUserInterface) AskForUntrustedFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool) {
	prompt := ui.FormatUntrustedFirstRunPrompt(bootDescriptor)

	prompt = userInterface.styleFirstRunPrompt(prompt)

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

func (userInterface *TerminalUserInterface) ShowLoader() {
	log.SetCallback(func(level logging.Level, message string) {
		if level <= logging.NOTICE {
			userInterface.SetStatus(message)
		}
	})
}

func (userInterface *TerminalUserInterface) HideLoader() {
	terminal := userInterface.terminal

	terminal.ResetStyle()
	terminal.ShowCursor()
	terminal.Clear()

	log.SetCallback(func(level logging.Level, message string) {})
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
	fmt.Printf("%v", userInterface.status)

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
	fmt.Printf(solidLine)

	terminal.EnableTextBold()
	terminal.PrintCenteredInRow(3, strings.ToUpper(userInterface.app))
	terminal.DisableTextBold()

	terminal.MoveCursor(5, 1)
	fmt.Printf(solidLine)
}
