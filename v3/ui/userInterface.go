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

import "github.com/giancosta86/moondeploy/v3/descriptors"

/*
UserInterface is the interface that must be implemented to plug a user interface,
based on any technology, into MoonDeploy's infrastructure
*/
type UserInterface interface {
	/*
		ShowError shows an error message
	*/
	ShowError(message string)

	/*
		AskForSecureFirstRun requests confirmation before running a secure app (ie,
		via HTTPS) for the first time.
		Must return true if the setup can proceed.
	*/
	AskForSecureFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool)

	/*
		AskForUntrustedFirstRun requests confirmation before running an untrusted
		app for the first time.
		Must return true if the setup can proceed.
	*/
	AskForUntrustedFirstRun(bootDescriptor descriptors.AppDescriptor) (canRun bool)

	/*
		SetApp sets the application name in the UI
	*/
	SetApp(app string)

	/*
		SetHeader sets the header describing a sequence of related activities
	*/
	SetHeader(header string)

	/*
		SetStatus logs a single activity to the UI
	*/
	SetStatus(status string)

	/*
		SetProgress sets the progress of the current activity, usually displayed
		via a progress bar. Must be in the range [0;1]
	*/
	SetProgress(progress float64)

	/*
		AskForDesktopShortcut asks the user whether a desktop shortcut should be created
		whenever an application has just been installed, Must return true if and
		only if the shortcut can be created
	*/
	AskForDesktopShortcut(referenceDescriptor descriptors.AppDescriptor) (canCreate bool)

	/*
		Show shows the user interface
	*/
	Show()

	/*
		Hide hides the user interface
	*/
	Hide()
}
