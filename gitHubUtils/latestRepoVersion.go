/*§
  ===========================================================================
  MoonDeploy
  ===========================================================================
  Copyright (C) 2015 Gianluca Costa
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

package gitHubUtils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/apps"
	"github.com/giancosta86/moondeploy/logging"
	"github.com/giancosta86/moondeploy/versioning"
)

var latestVersionUrlRegex = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/releases/latest/?`)
var tagRegex = regexp.MustCompile(`^\D*(\d.*)`)

var apiLatestVersioURLTemplate = "https://api.github.com/repos/%v/%v/releases/latest"

type assetInfo struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type apiLatestVersionResponse struct {
	TagName string      `json:"tag_name"`
	Assets  []assetInfo `json:"assets"`
}

type GitHubRemoteDescriptorInfo struct {
	Version       *versioning.Version
	DescriptorURL *url.URL
}

func GetLatestRemoteDescriptorInfo(baseUrl *url.URL) *GitHubRemoteDescriptorInfo {
	logging.Info("Checking if the Base URL matches GitHub's /latest release URL pattern...")

	gitHubProjectParams := latestVersionUrlRegex.FindStringSubmatch(baseUrl.String())
	if gitHubProjectParams == nil {
		logging.Info("The URL does not match")
		return nil
	}
	logging.Notice("The URL matches")

	gitHubUser := gitHubProjectParams[1]
	gitHubRepo := gitHubProjectParams[2]

	apiLatestVersioURL, err := url.Parse(fmt.Sprintf(
		apiLatestVersioURLTemplate,
		gitHubUser,
		gitHubRepo))
	if err != nil {
		logging.Warning(err.Error())
		return nil
	}

	logging.Info("Calling GitHub's API, at '%v'...", apiLatestVersioURL)

	responseBytes, err := caravel.RetrieveFromURL(apiLatestVersioURL)
	if err != nil {
		logging.Warning(err.Error())
		return nil
	}
	logging.Notice("API returned OK")

	logging.Info("Deserializing the API response...")
	var latestVersionResponse apiLatestVersionResponse
	err = json.Unmarshal(responseBytes, &latestVersionResponse)
	if err != nil {
		logging.Warning(err.Error())
		return nil
	}
	logging.Notice("Response correctly deserialized: %#v", latestVersionResponse)

	logging.Info("Now processing the response fields...")

	result := GitHubRemoteDescriptorInfo{}

	for _, asset := range latestVersionResponse.Assets {
		if asset.Name == apps.DescriptorFileName {
			result.DescriptorURL, err = url.Parse(asset.BrowserDownloadURL)
			if err != nil {
				logging.Warning(err.Error())
				return nil
			}
			break
		}
	}

	if result.DescriptorURL == nil {
		logging.Warning("The app descriptor could not be found as an asset of the latest release")
		return nil
	}

	tagComponents := tagRegex.FindStringSubmatch(latestVersionResponse.TagName)
	if tagComponents == nil {
		logging.Warning("GitHub's release tag must be in the format: <string or empty><VERSION>, not '%v'", latestVersionResponse.TagName)
		return nil
	}

	result.Version, err = versioning.ParseVersion(tagComponents[1])
	if err != nil {
		logging.Warning(err.Error())
		return nil
	}

	logging.Notice("Response fields correctly processed")

	return &result
}
