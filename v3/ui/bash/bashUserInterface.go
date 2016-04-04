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

package bash

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/logging"
	"github.com/giancosta86/moondeploy/v3/ui"
)

type BashUserInterface struct {
	app      string
	header   string
	status   string
	progress float64

	rows      int
	columns   int
	emptyRow  string
	strongRow string
}

func NewBashUserInterface() *BashUserInterface {
	rows, columns := getTerminalSize()

	logging.Notice("Rows: %v", rows)
	logging.Notice("Columns: %v", columns)

	emptyRow := strings.Repeat(" ", columns)
	strongRow := strings.Repeat("-", columns)

	return &BashUserInterface{
		rows:      rows,
		columns:   columns,
		emptyRow:  emptyRow,
		strongRow: strongRow,
	}
}

func getTerminalSize() (rows int, columns int) {
	const defaultRows = 25
	const defaultColumns = 80

	sizeCommand := exec.Command("stty", "size")
	sizeCommand.Stdin = os.Stdin

	sizeOutputString, err := caravel.GetCommandOutputString(sizeCommand)
	if err != nil {
		logging.Warning("Could not run stty: %v", err.Error())
	} else {
		sizeOutputString = strings.TrimSpace(sizeOutputString)
		logging.Notice("Retrieved size values: %v", sizeOutputString)
	}

	_, err = fmt.Sscanf(sizeOutputString, "%d %d", &rows, &columns)
	if err != nil {
		logging.Warning("Cannot parse the size string: %v", err.Error())

		return defaultRows, defaultColumns
	}

	return rows, columns
}

func (userInterface *BashUserInterface) ShowError(message string) {
	userInterface.HideLoader()

	logging.Error(message)
}

func (userInterface *BashUserInterface) askQuestion(prompt string) (yesResponse bool) {
	for {
		userInterface.clearScreen()
		userInterface.setupColors()
		userInterface.drawEmptyBackground()
		userInterface.drawTitle()

		userInterface.moveCursorTo(8, 0)
		fmt.Printf("%v [Y/N]: ", prompt)

		reader := bufio.NewReader(os.Stdin)
		userInput, err := reader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		userInput = strings.TrimSpace(userInput)

		switch userInput {
		case "Y":
			return true
		case "N":
			return false
		}

	}
}

func (userInterface *BashUserInterface) styleFirstRunPrompt(prompt string) string {
	prompt = strings.Replace(prompt, "\nTitle:", "\n\033[1m   Title:\033[21m", -1)
	prompt = strings.Replace(prompt, "\nPublisher:", "\n\033[1m   Publisher:\033[21m", -1)
	prompt = strings.Replace(prompt, "\nAddress:", "\n\033[1m   Address:\033[21m", -1)

	return prompt
}

func (userInterface *BashUserInterface) AskForSecureFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool) {
	prompt := ui.FormatSecureFirstRunPrompt(bootDescriptor)

	prompt = userInterface.styleFirstRunPrompt(prompt)

	return userInterface.askQuestion(prompt)
}

func (userInterface *BashUserInterface) AskForUntrustedFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool) {
	prompt := ui.FormatUntrustedFirstRunPrompt(bootDescriptor)

	prompt = userInterface.styleFirstRunPrompt(prompt)

	prompt = strings.Replace(prompt, "\nWARNING:", "\n\033[1mWARNING:\033[21m", -1)

	return userInterface.askQuestion(prompt)
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
	prompt := ui.FormatDesktopShortcutPrompt(referenceDescriptor)
	return userInterface.askQuestion(prompt)
}

func (userInterface *BashUserInterface) ShowLoader() {
	logging.SetOutputEnabled(false)
	logging.SetCallback(func(message string) {
		userInterface.SetStatus(message)
	})
}

func (userInterface *BashUserInterface) HideLoader() {
	userInterface.resetTerminal()
	userInterface.clearScreen()

	logging.SetCallback(func(message string) {})

	logging.SetOutputEnabled(true)
}

func (userInterface *BashUserInterface) setupColors() {
	//Background
	fmt.Print("\033[48;5;195m")

	//Foreground
	fmt.Print("\033[38;5;22m")
}

func (userInterface *BashUserInterface) clearScreen() {
	fmt.Print("\033c")

	userInterface.moveCursorTo(0, 0)
}

func (userInterface *BashUserInterface) moveCursorTo(rowIndex int, columnIndex int) {
	fmt.Printf("\033[%d;%dH", rowIndex, columnIndex)
}

func (userInterface *BashUserInterface) resetTerminal() {
	fmt.Print("\033[0m")
}

func (userInterface *BashUserInterface) enableHiddenText() {
	fmt.Print("\033[8m")
}

func (userInterface *BashUserInterface) enableBold() {
	fmt.Print("\033[1m")
}

func (userInterface *BashUserInterface) disableBold() {
	fmt.Print("\033[21m")
}

func (userInterface *BashUserInterface) enableUnderlined() {
	fmt.Print("\033[4m")
}

func (userInterface *BashUserInterface) disableUnderlined() {
	fmt.Print("\033[24m")
}

func (userInterface *BashUserInterface) redraw() {
	userInterface.clearScreen()
	userInterface.setupColors()

	userInterface.drawEmptyBackground()

	userInterface.drawTitle()

	userInterface.enableBold()
	userInterface.printCentered(8, userInterface.header)
	userInterface.disableBold()

	userInterface.moveCursorTo(12, 2)
	fmt.Printf("%v", userInterface.status)

	if userInterface.progress > 0 {
		userInterface.drawProgressBar(16, 2, userInterface.columns-20, userInterface.progress)
	}

	userInterface.enableHiddenText()

	//TODO: DEL THIS
	time.Sleep(500 * time.Millisecond)
}

func (userInterface *BashUserInterface) drawEmptyBackground() {
	for i := 1; i <= userInterface.rows-1; i++ {
		fmt.Println(userInterface.emptyRow)
	}
	fmt.Print(userInterface.emptyRow)
}

func (userInterface *BashUserInterface) drawTitle() {
	if userInterface.app == "" {
		return
	}

	userInterface.moveCursorTo(0, 0)
	fmt.Printf(userInterface.strongRow)

	userInterface.enableBold()
	userInterface.printCentered(3, strings.ToUpper(userInterface.app))
	userInterface.disableBold()

	userInterface.moveCursorTo(5, 0)
	fmt.Printf(userInterface.strongRow)

}

func (userInterface *BashUserInterface) printCentered(rowIndex int, text string) {
	textLength := utf8.RuneCountInString(text)
	textColumnIndex := (userInterface.columns - textLength) / 2

	userInterface.moveCursorTo(rowIndex, textColumnIndex)
	fmt.Printf("%v", text)
}

func (userInterface *BashUserInterface) drawProgressBar(rowIndex int, startColIndex int, maxNumberOfTicks int, fractionalValue float64) {
	const delimiter = "ǁ"
	const tick = "="
	const space = " "

	numberOfTicks := int(math.Ceil(float64(maxNumberOfTicks) * fractionalValue))
	percentage := 100 * fractionalValue

	userInterface.moveCursorTo(rowIndex, startColIndex)
	fmt.Printf("%v%v%v%v  %.2f%%",
		delimiter,
		strings.Repeat(tick, numberOfTicks),
		strings.Repeat(space, maxNumberOfTicks-numberOfTicks),
		delimiter,
		percentage)
}
