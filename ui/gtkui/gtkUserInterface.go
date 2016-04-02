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

package gtkui

import (
	"fmt"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"

	"github.com/giancosta86/moondeploy"
	"github.com/giancosta86/moondeploy/apps"
)

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
		builder, err := gtk.BuilderNew()
		if err != nil {
			panic(err)
		}

		gladeDescriptorPath := filepath.Join(
			filepath.Dir(moondeploy.Executable),
			"moondeploy.glade")

		err = builder.AddFromFile(gladeDescriptorPath)
		if err != nil {
			panic(err)
		}

		windowObject, err := builder.GetObject("mainWindow")
		if err != nil {
			panic(err)
		}
		window := windowObject.(*gtk.Window)
		userInterface.window = window

		window.SetTitle(moondeploy.Title)
		window.SetIconFromFile(moondeploy.IconPathAsPng)

		window.Connect("destroy", func() {
			window.Destroy()
			userInterface.window = nil
			gtk.MainQuit()
			userInterface.closedByUser = true
		})

		appLabelObject, err := builder.GetObject("appLabel")
		if err != nil {
			panic(err)
		}
		userInterface.appLabel = appLabelObject.(*gtk.Label)

		headerLabelObject, err := builder.GetObject("headerLabel")
		if err != nil {
			panic(err)
		}
		userInterface.headerLabel = headerLabelObject.(*gtk.Label)

		statusLabelObject, err := builder.GetObject("statusLabel")
		if err != nil {
			panic(err)
		}
		userInterface.statusLabel = statusLabelObject.(*gtk.Label)

		progressBarObject, err := builder.GetObject("progressBar")
		if err != nil {
			panic(err)
		}
		userInterface.progressBar = progressBarObject.(*gtk.ProgressBar)

		return nil
	})

	return userInterface, nil
}

func (userInterface *GtkUserInterface) IsClosedByUser() bool {
	return userInterface.closedByUser
}

func (userInterface *GtkUserInterface) showBasicMessageDialog(messageType gtk.MessageType, message string) {
	runOnUIThreadAndWait(func() interface{} {
		dialog := gtk.MessageDialogNew(userInterface.window, gtk.DIALOG_MODAL, messageType, gtk.BUTTONS_OK, message)
		defer dialog.Destroy()

		dialog.SetTitle(moondeploy.Title)
		dialog.Run()

		return nil
	})
}

func (userInterface *GtkUserInterface) showYesNoDialog(messageType gtk.MessageType, message string) bool {
	result := runOnUIThreadAndWait(func() interface{} {
		dialog := gtk.MessageDialogNew(userInterface.window, gtk.DIALOG_MODAL, messageType, gtk.BUTTONS_YES_NO, message)
		defer dialog.Destroy()

		dialog.SetTitle(moondeploy.Title)
		return (dialog.Run() == int(gtk.RESPONSE_YES))
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

func formatBasicFirstRunPrompt(bootDescriptor apps.AppDescriptor) (basicFirstRunPrompt string) {
	const basicFirstRunTemplate = "You are running an application for the first time." +
		"\n\n\nTitle:   %v" +
		"\n\nPublisher:   %v" +
		"\n\nAddress:   %v\n\n\nDo you wish to proceed?"

	return fmt.Sprintf(basicFirstRunTemplate,

		bootDescriptor.GetTitle(),
		bootDescriptor.GetPublisher(),
		bootDescriptor.GetDeclaredBaseURL())
}

func (userInterface *GtkUserInterface) AskForSecureFirstRun(bootDescriptor apps.AppDescriptor) (canRun bool) {
	return userInterface.askYesNo(formatBasicFirstRunPrompt(bootDescriptor))
}

func (userInterface *GtkUserInterface) AskForUntrustedFirstRun(bootDescriptor apps.AppDescriptor) (canRun bool) {
	return userInterface.askWarningYesNo(
		formatBasicFirstRunPrompt(bootDescriptor) + untrustedWarning)
}

func (userInterface *GtkUserInterface) AskForDesktopShortcut(referenceDescriptor apps.AppDescriptor) (canCreate bool) {
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
