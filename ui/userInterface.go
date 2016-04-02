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

import "github.com/giancosta86/moondeploy/apps"

/*
UserInterface is the interface that must be implemented to plug a user interface,
base on any technology, into MoonDeploy's algorithm
*/
type UserInterface interface {
	ShowError(message string)

	AskForSecureFirstRun(bootDescriptor apps.AppDescriptor) (canRun bool)
	AskForUntrustedFirstRun(bootDescriptor apps.AppDescriptor) (canRun bool)

	SetApp(app string)
	SetHeader(header string)
	SetStatus(status string)
	SetProgress(progress float64)

	AskForDesktopShortcut(referenceDescriptor apps.AppDescriptor) (canCreate bool)

	HideLoader()
}
