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
	"github.com/giancosta86/caravel/terminals"

	"github.com/giancosta86/moondeploy/v3/descriptors"
	"github.com/giancosta86/moondeploy/v3/engine"
	"github.com/giancosta86/moondeploy/v3/launchers"
	"github.com/giancosta86/moondeploy/v3/log"
	"github.com/giancosta86/moondeploy/v3/ui/termui"
)

func StartGUI(launcher launchers.Launcher, bootDescriptorPath string) (err error) {
	bootDescriptor, err := descriptors.NewAppDescriptorFromPath(bootDescriptorPath)
	if err != nil {
		return err
	}

	bashTerminal := terminals.NewBashTerminal()

	userInterface := termui.NewTerminalUserInterface(launcher, bashTerminal)

	userInterface.Show()

	result := engine.Run(launcher, userInterface, bootDescriptor)

	userInterface.Hide()

	log.Notice("OK")

	return result
}
