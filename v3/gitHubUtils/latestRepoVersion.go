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

package gitHubUtils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy/v3/log"
	"github.com/giancosta86/moondeploy/v3/versioning"
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

func GetLatestRemoteDescriptorInfo(baseUrl *url.URL, descriptorFileName string) *GitHubRemoteDescriptorInfo {
	log.Info("Checking if the Base URL matches GitHub's /latest release URL pattern...")

	gitHubProjectParams := latestVersionUrlRegex.FindStringSubmatch(baseUrl.String())
	if gitHubProjectParams == nil {
		log.Info("The URL does not match a GitHub pattern")
		return nil
	}
	log.Notice("The URL matches a GitHub 'latest' pattern")

	gitHubUser := gitHubProjectParams[1]
	gitHubRepo := gitHubProjectParams[2]

	apiLatestVersioURL, err := url.Parse(fmt.Sprintf(
		apiLatestVersioURLTemplate,
		gitHubUser,
		gitHubRepo))
	if err != nil {
		log.Warning(err.Error())
		return nil
	}

	log.Info("Calling GitHub's API, at '%v'...", apiLatestVersioURL)

	responseBytes, err := caravel.RetrieveFromURL(apiLatestVersioURL)
	if err != nil {
		log.Warning(err.Error())
		return nil
	}
	log.Notice("API returned OK")

	log.Info("Deserializing the API response...")
	var latestVersionResponse apiLatestVersionResponse
	err = json.Unmarshal(responseBytes, &latestVersionResponse)
	if err != nil {
		log.Warning(err.Error())
		return nil
	}
	log.Info("Response correctly deserialized: %#v", latestVersionResponse)

	log.Info("Now processing the response fields...")

	result := GitHubRemoteDescriptorInfo{}

	for _, asset := range latestVersionResponse.Assets {
		if asset.Name == descriptorFileName {
			result.DescriptorURL, err = url.Parse(asset.BrowserDownloadURL)
			if err != nil {
				log.Warning(err.Error())
				return nil
			}
			break
		}
	}

	if result.DescriptorURL == nil {
		log.Warning("The app descriptor could not be found as an asset of the latest release")
		return nil
	}

	tagComponents := tagRegex.FindStringSubmatch(latestVersionResponse.TagName)
	if tagComponents == nil {
		log.Warning("GitHub's release tag must be in the format: <string or empty><VERSION>, not '%v'", latestVersionResponse.TagName)
		return nil
	}

	result.Version, err = versioning.ParseVersion(tagComponents[1])
	if err != nil {
		log.Warning(err.Error())
		return nil
	}

	log.Notice("Response fields correctly processed")

	return &result
}
