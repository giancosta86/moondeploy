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

//TODO: restore these lines and move the lines below to the Darwin implementation
/*package verbs

import (
	"github.com/giancosta86/moondeploy/v3/custom"
	"github.com/giancosta86/moondeploy/v3/moonclient/gtkLauncher"
)

func StartGUI(bootDescriptorPath string, settings *custom.Settings) (err error) {
	return gtkLauncher.StartGUI(bootDescriptorPath, settings)
}
*/

package verbs

import (
	"github.com/giancosta86/moondeploy/v3/custom"
	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/engine"
	"github.com/giancosta86/moondeploy/v3/logging"
	"github.com/giancosta86/moondeploy/v3/ui"

	goLogging "github.com/op/go-logging"
)

func StartGUI(bootDescriptorPath string, settings *custom.Settings) (err error) {
	logging.SetLevel(goLogging.CRITICAL)

	bootDescriptor, err := descriptors.NewAppDescriptorFromPath(bootDescriptorPath)
	if err != nil {
		return err
	}

	userInterface := ui.NewBashUserInterface()

	return engine.Run(bootDescriptor, settings, userInterface)
}
