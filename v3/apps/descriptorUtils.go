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

package apps

import (
	"fmt"
	"net/url"
	"path"
)

func getRelativeFileURL(descriptor AppDescriptor, relativePath string) (*url.URL, error) {
	if path.IsAbs(relativePath) {
		return nil, fmt.Errorf("Absolute paths are not allowed: '%v'", relativePath)
	}

	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return nil, err
	}

	return descriptor.GetActualBaseURL().ResolveReference(relativeURL), nil
}
