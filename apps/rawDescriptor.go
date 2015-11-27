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

package apps

import (
	"fmt"
	"net/url"

	"github.com/giancosta86/moondeploy/logging"
	"github.com/giancosta86/moondeploy/versioning"
)

type rawAppDescriptor struct {
	DescriptorVersion string

	Name            string
	Version         string
	SkipUpdateCheck bool

	Description string
	IconPath    map[string]string

	BaseURL   string
	Publisher string

	PackageVersions map[string]string

	CommandLine map[string][]string

	SkipPackageLevels int
}

func (rawDescriptor *rawAppDescriptor) toFull() (descriptor *AppDescriptor, err error) {
	err = rawDescriptor.rawValidate()
	if err != nil {
		return nil, err
	}

	descriptor = &AppDescriptor{}

	descriptor.DescriptorVersion, err = versioning.ParseVersion(rawDescriptor.DescriptorVersion)
	if err != nil {
		return nil, err
	}

	descriptor.Name = rawDescriptor.Name

	descriptor.Version, err = versioning.ParseVersion(rawDescriptor.Version)
	if err != nil {
		return nil, err
	}

	descriptor.SkipUpdateCheck = rawDescriptor.SkipUpdateCheck

	descriptor.Description = rawDescriptor.Description
	descriptor.IconPath = rawDescriptor.IconPath

	descriptor.BaseURL, err = url.Parse(rawDescriptor.BaseURL)
	if err != nil {
		return nil, err
	}

	descriptor.Publisher = rawDescriptor.Publisher

	descriptor.PackageVersions = make(map[string]*versioning.Version)
	for packageName, packageVersionString := range rawDescriptor.PackageVersions {
		if packageVersionString != "" {
			descriptor.PackageVersions[packageName], err = versioning.ParseVersion(packageVersionString)
			if err != nil {
				logging.Warning("Invalid version string for package '%v': '%v'; assuming it is missing",
					packageName,
					packageVersionString)
				descriptor.PackageVersions[packageName] = nil
			}
		} else {
			descriptor.PackageVersions[packageName] = nil
		}
	}

	descriptor.CommandLine = rawDescriptor.CommandLine

	descriptor.SkipPackageLevels = rawDescriptor.SkipPackageLevels

	return descriptor, nil
}

func (rawDescriptor *rawAppDescriptor) rawValidate() (err error) {
	if rawDescriptor.DescriptorVersion == "" {
		return fmt.Errorf("Missing descriptor version")
	}

	if rawDescriptor.Version == "" {
		return fmt.Errorf("Missing version")
	}

	if rawDescriptor.BaseURL == "" {
		return fmt.Errorf("Missing base URL")
	}

	if rawDescriptor.IconPath == nil {
		rawDescriptor.IconPath = make(map[string]string)
	}

	if rawDescriptor.BaseURL[len(rawDescriptor.BaseURL)-1] != '/' {
		rawDescriptor.BaseURL = rawDescriptor.BaseURL + "/"
	}

	if rawDescriptor.PackageVersions == nil {
		rawDescriptor.PackageVersions = make(map[string]string)
	}

	if rawDescriptor.CommandLine == nil {
		rawDescriptor.CommandLine = make(map[string][]string)
	}

	return nil
}
