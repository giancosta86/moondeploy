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

package verbs

import (
	"os"
	"strconv"

	"github.com/giancosta86/moondeploy/v3/logging"
	"github.com/giancosta86/moondeploy/v3/server"
)

func DoServe() (err error) {
	if len(os.Args) < 4 {
		return &InvalidCommandLineArguments{}
	}

	portString := os.Args[2]
	sourceDir := os.Args[3]

	port, err := strconv.Atoi(portString)
	if err != nil {
		return err
	}

	logging.Info("Activating server on port %v for dir: '%v'...", port, sourceDir)

	return server.ServeDirectory(sourceDir, port)
}
