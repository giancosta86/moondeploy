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

package gtkui

import (
	"fmt"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"

	"github.com/giancosta86/moondeploy"
	"github.com/giancosta86/moondeploy/apps"
)

const basicFirstRunTemplate = "You are running an application for the first time." +
	"\n\n\nTitle:   %v" +
	"\n\nPublisher:   %v" +
	"\n\nAddress:   %v\n\n\nDo you wish to proceed?"

const untrustedWarning = "\n\n\nWARNING: the provided address is insecure, so " +
	"the integrity of the application files might be compromised by " +
	"third parties during the download process."

type GtkUserInterface struct {
	window *gtk.Window

	appLabel *gtk.Label

	headerLabel *gtk.Label
	statusLabel *gtk.Label
	progressBar *gtk.ProgressBar

	closedByUser bool
}

func NewGtkUserInterface() (userInterface *GtkUserInterface, err error) {
	userInterface = &GtkUserInterface{}

	runOnUIThreadAndWait(func() interface{} {
		window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
		userInterface.window = window

		window.SetTitle(moondeploy.Title)
		window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)
		window.SetTypeHint(gdk.WINDOW_TYPE_HINT_DIALOG)
		window.SetIconFromFile(moondeploy.IconPathAsPng)
		window.SetSizeRequest(800, 350)

		window.Connect("destroy", func(ctx *glib.CallbackContext) {
			window.Destroy()
			userInterface.window = nil
			gtk.MainQuit()
			userInterface.closedByUser = true
		})

		mainBox := gtk.NewVBox(false, 10)
		mainAlignment := gtk.NewAlignment(0, 0, 1, 1)
		mainAlignment.SetPadding(40, 40, 40, 40)
		mainAlignment.Add(mainBox)
		window.Add(mainAlignment)

		appLabel := gtk.NewLabel("")
		userInterface.appLabel = appLabel
		appLabel.ModifyFontEasy("bold 24")
		mainBox.PackStart(appLabel, false, false, 15)

		appBox := gtk.NewHBox(false, 40)
		mainBox.PackStart(appBox, true, true, 20)

		spinner := gtk.NewSpinner()
		spinner.SetSizeRequest(32, 32)
		spinner.Start()
		appBox.PackStart(spinner, false, false, 0)

		infoBox := gtk.NewVBox(false, 20)
		appBox.PackStart(infoBox, true, true, 0)

		headerLabel := gtk.NewLabel("")
		userInterface.headerLabel = headerLabel
		headerLabel.SetAlignment(0, 0.5)
		headerLabel.ModifyFontEasy("bold")
		infoBox.PackStart(headerLabel, true, true, 0)

		statusLabel := gtk.NewLabel("")
		userInterface.statusLabel = statusLabel
		statusLabel.SetAlignment(0, 0.5)
		statusLabel.SetLineWrap(true)
		infoBox.PackStart(statusLabel, true, true, 0)

		progressBar := gtk.NewProgressBar()
		userInterface.progressBar = progressBar
		progressBar.SetSizeRequest(-1, 20)
		infoBox.PackStart(progressBar, true, false, 0)

		return nil
	})

	return userInterface, nil
}

func (userInterface *GtkUserInterface) IsClosedByUser() bool {
	return userInterface.closedByUser
}

func (userInterface *GtkUserInterface) showBasicMessageDialog(messageType gtk.MessageType, message string) {
	runOnUIThreadAndWait(func() interface{} {
		dialog := gtk.NewMessageDialog(userInterface.window, gtk.DIALOG_MODAL, messageType, gtk.BUTTONS_OK, message)
		defer dialog.Destroy()

		dialog.SetTitle(moondeploy.Title)
		dialog.Run()

		return nil
	})
}

func (userInterface *GtkUserInterface) showYesNoDialog(messageType gtk.MessageType, message string) bool {
	result := runOnUIThreadAndWait(func() interface{} {
		dialog := gtk.NewMessageDialog(userInterface.window, gtk.DIALOG_MODAL, messageType, gtk.BUTTONS_YES_NO, message)
		defer dialog.Destroy()

		dialog.SetTitle(moondeploy.Title)
		return (dialog.Run() == gtk.RESPONSE_YES)
	})

	return result.(bool)
}

func (userInterface *GtkUserInterface) ShowError(message string) {
	userInterface.showBasicMessageDialog(gtk.MESSAGE_ERROR, message)
}

func (userInterface *GtkUserInterface) askYesNo(message string) bool {
	return userInterface.showYesNoDialog(gtk.MESSAGE_QUESTION, message)
}

func (userInterface *GtkUserInterface) askWarningYesNo(message string) bool {
	return userInterface.showYesNoDialog(gtk.MESSAGE_WARNING, message)
}

func formatBasicFirstRunPrompt(bootDescriptor *apps.AppDescriptor) (basicFirstRunPrompt string) {
	return fmt.Sprintf(basicFirstRunTemplate,

		bootDescriptor.GetTitle(),
		bootDescriptor.Publisher,
		bootDescriptor.BaseURL)
}

func (userInterface *GtkUserInterface) AskForSecureFirstRun(bootDescriptor *apps.AppDescriptor) (canRun bool) {
	return userInterface.askYesNo(formatBasicFirstRunPrompt(bootDescriptor))
}

func (userInterface *GtkUserInterface) AskForUntrustedFirstRun(bootDescriptor *apps.AppDescriptor) (canRun bool) {
	return userInterface.askWarningYesNo(
		formatBasicFirstRunPrompt(bootDescriptor) + untrustedWarning)
}

func (userInterface *GtkUserInterface) AskForDesktopShortcut(referenceDescriptor *apps.AppDescriptor) (canCreate bool) {
	return userInterface.askYesNo("Would you like to create a desktop shortcut for the application?")
}

func (userInterface *GtkUserInterface) SetApp(app string) {
	runOnUIThreadAndWait(func() interface{} {
		userInterface.window.SetTitle(fmt.Sprintf("%v - %v", moondeploy.Name, app))
		userInterface.appLabel.SetText(app)
		return nil
	})
}

func (userInterface *GtkUserInterface) SetHeader(header string) {
	runOnUIThreadAndWait(func() interface{} {
		userInterface.headerLabel.SetText(header)
		return nil
	})
}

func (userInterface *GtkUserInterface) SetStatus(status string) {
	runOnUIThreadAndWait(func() interface{} {
		userInterface.statusLabel.SetText(status)
		return nil
	})
}

func (userInterface *GtkUserInterface) SetProgress(progress float64) {
	runOnUIThreadAndWait(func() interface{} {
		userInterface.progressBar.SetFraction(progress)

		if progress > 0 {
			userInterface.progressBar.SetVisible(true)
		} else {
			userInterface.progressBar.SetVisible(false)
		}

		return nil
	})
}

func (userInterface *GtkUserInterface) ShowLoader() {
	runOnUIThreadAndWait(func() interface{} {
		userInterface.window.ShowAll()
		userInterface.progressBar.Hide()

		return nil
	})
}

func (userInterface *GtkUserInterface) HideLoader() {
	runOnUIThreadAndWait(func() interface{} {
		userInterface.window.Hide()
		return nil
	})
}

func InitGTK() {
	gtk.Init(nil)
}
