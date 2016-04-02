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
	"net/url"

	"github.com/giancosta86/moondeploy/v3/gitHubUtils"
	"github.com/giancosta86/moondeploy/v3/logging"
)

type actualBaseURLSearchStrategy func(AppDescriptor) *url.URL

var actualBaseURLCache = make(map[string]*url.URL)

var actualBaseURLSearchStrategies = []actualBaseURLSearchStrategy{
	lookForActualURLInCache,
	lookForActualURLOnGitHub}

func getActualBaseURL(descriptor AppDescriptor) *url.URL {
	var actualBaseURL *url.URL

	for _, searchStrategy := range actualBaseURLSearchStrategies {
		actualBaseURL = searchStrategy(descriptor)

		if actualBaseURL != nil {
			logging.Notice("The actual base URL has been found by a search strategy!")
			break
		}
	}

	if actualBaseURL == nil {
		logging.Info("The actual base URL just matches the actual base URL")
		actualBaseURL = descriptor.GetDeclaredBaseURL()
	}

	actualBaseURLCache[descriptor.GetDeclaredBaseURL().String()] = actualBaseURL
	logging.Info("Actual base URL '%v' put into the cache", actualBaseURL)

	return actualBaseURL
}

func lookForActualURLInCache(descriptor AppDescriptor) *url.URL {
	logging.Info("Checking if the Base URL is a key of the actual Base URL cache...")

	cachedActualURL, _ := actualBaseURLCache[descriptor.GetDeclaredBaseURL().String()]

	if cachedActualURL != nil {
		logging.Notice("Actual URL found in the cache! --> '%v'", cachedActualURL)

		return cachedActualURL
	}

	logging.Info("Actual URL not in the cache")
	return nil
}

func lookForActualURLOnGitHub(descriptor AppDescriptor) *url.URL {
	logging.Info("Checking if the Base URL points to the *latest* release of a GitHub repo...")
	gitHubLatestRemoteDescriptorInfo := gitHubUtils.GetLatestRemoteDescriptorInfo(
		descriptor.GetDeclaredBaseURL(),
		descriptor.GetDescriptorFileName())

	if gitHubLatestRemoteDescriptorInfo != nil {
		logging.Notice("The given base URL actually references version '%v', whose descriptor is at URL: '%v'",
			gitHubLatestRemoteDescriptorInfo.Version,
			gitHubLatestRemoteDescriptorInfo.DescriptorURL)

		parentDirURL, err := url.Parse(".")
		if err != nil {
			panic(err)
		}

		actualBaseURL := gitHubLatestRemoteDescriptorInfo.DescriptorURL.ResolveReference(parentDirURL)

		logging.Notice("The actual base URL returned by GitHub is: '%v'", actualBaseURL)
		return actualBaseURL
	}

	return nil
}
