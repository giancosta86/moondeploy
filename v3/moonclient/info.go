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

package moonclient

import (
	"net/url"
	"path/filepath"
	"runtime"

	"github.com/kardianos/osext"

	"github.com/giancosta86/moondeploy"
)

const Name = "MoonDeploy"

var Title = Name + " " + moondeploy.Version

var WebsiteURL *url.URL

var Executable string
var Dir string

var IconPathAsIco string
var IconPathAsPng string

func GetIconPath() string {
	if runtime.GOOS == "windows" {
		return IconPathAsIco
	}

	return IconPathAsPng
}

func init() {
	var err error

	WebsiteURL, err = url.Parse("https://github.com/giancosta86/moondeploy")
	if err != nil {
		panic(err)
	}

	Executable, err = osext.Executable()
	if err != nil {
		panic(err)
	}

	Dir, err = osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

	IconPathAsIco = filepath.Join(Dir, "moondeploy.ico")
	IconPathAsPng = filepath.Join(Dir, "moondeploy.png")
}
