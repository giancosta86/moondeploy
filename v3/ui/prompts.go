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
)

func FormatSecureFirstRunPrompt(bootDescriptor descriptors.AppDescriptor) (basicFirstRunPrompt string) {
	const basicFirstRunTemplate = "You are running an application for the first time." +
		"\n\n\nTitle:   %v" +
		"\n\nPublisher:   %v" +
		"\n\nAddress:   %v\n\n\nDo you wish to proceed?"

	return fmt.Sprintf(basicFirstRunTemplate,

		bootDescriptor.GetTitle(),
		bootDescriptor.GetPublisher(),
		bootDescriptor.GetDeclaredBaseURL())
}

func FormatUntrustedFirstRunPrompt(bootDescriptor descriptors.AppDescriptor) (basicFirstRunPrompt string) {
	return FormatSecureFirstRunPrompt(bootDescriptor) + untrustedWarning
}

func FormatDesktopShortcutPrompt(referenceDescriptor descriptors.AppDescriptor) (prompt string) {
	return "Would you like to create a desktop shortcut for the application?"
}

const untrustedWarning = "\n\n\nWARNING: the provided address is insecure, so " +
	"the integrity of the application files might be compromised by " +
	"third parties during the download process."
