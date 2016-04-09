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

package descriptors

import (
	"fmt"
	"net/url"
	"path"

	"github.com/giancosta86/moondeploy/v3/versioning"
)

func parsePackageVersions(packageVersionsStringMap map[string]string) (result map[string]*versioning.Version, err error) {
	result = make(map[string]*versioning.Version)

	if packageVersionsStringMap == nil {
		return result, nil
	}

	for packageName, packageVersionString := range packageVersionsStringMap {
		if packageVersionString != "" {
			result[packageName], err = versioning.ParseVersion(packageVersionString)
			if err != nil {
				return nil, fmt.Errorf("Invalid version string for package '%v': '%v'",
					packageName,
					packageVersionString)
			}
		} else {
			result[packageName] = nil
		}
	}

	return result, nil
}

func getRemoteFileURL(descriptor AppDescriptor, relativePath string) (*url.URL, error) {
	if path.IsAbs(relativePath) {
		return nil, fmt.Errorf("Absolute paths are not allowed: '%v'", relativePath)
	}

	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return nil, err
	}

	return descriptor.GetActualBaseURL().ResolveReference(relativeURL), nil
}

func ensureTrailingSlash(path string) string {
	if path[len(path)-1] != '/' {
		return path + "/"
	}

	return path
}
