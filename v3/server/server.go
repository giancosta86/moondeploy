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

package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/giancosta86/moondeploy/v3/log"
)

/*
ServeDirectory activates a basic static web server on the given port, serving
from the given source dir.
A client can stop it by accessing its "/moondeploy.quit" path.
*/
func ServeDirectory(sourceDir string, port int) (err error) {
	fileServer := http.FileServer(http.Dir(sourceDir))
	http.Handle("/", fileServer)

	http.HandleFunc("/moondeploy.quit", func(http.ResponseWriter, *http.Request) {
		log.Notice("Now quitting")
		os.Exit(0)
	})

	portString := ":" + strconv.Itoa(port)

	return http.ListenAndServe(portString, nil)
}
