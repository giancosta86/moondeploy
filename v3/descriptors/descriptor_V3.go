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
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"

	"github.com/giancosta86/moondeploy/v3/versioning"
)

type appDescriptorV3 struct {
	DescriptorVersion string

	BaseURL            string
	DescriptorFileName string

	Name        string
	Version     string
	Publisher   string
	Description string

	SkipPackageLevels int
	SkipUpdateCheck   bool

	SupportedOS []string

	osSettingsStruct

	OS map[string]osSettingsStruct

	//
	//Computed fields
	//
	descriptorVersion *versioning.Version

	declaredBaseURL    *url.URL
	actualBaseURL      *url.URL
	descriptorFileName string

	name        string
	appVersion  *versioning.Version
	publisher   string
	description string

	skipPackageLevels int
	skipUpdateCheck   bool

	supportedSystems []string

	packageVersions map[string]*versioning.Version
	commandLine     []string
	iconPath        string
}

type osSettingsStruct struct {
	Packages    map[string]string
	CommandLine []string
	IconPath    string
}

func (descriptor *appDescriptorV3) GetDescriptorVersion() *versioning.Version {
	return descriptor.descriptorVersion
}

func (descriptor *appDescriptorV3) GetDeclaredBaseURL() *url.URL {
	return descriptor.declaredBaseURL
}

func (descriptor *appDescriptorV3) GetActualBaseURL() *url.URL {
	return descriptor.actualBaseURL
}

func (descriptor *appDescriptorV3) GetDescriptorFileName() string {
	return descriptor.descriptorFileName
}

func (descriptor *appDescriptorV3) GetName() string {
	return descriptor.name
}

func (descriptor *appDescriptorV3) GetAppVersion() *versioning.Version {
	return descriptor.appVersion
}

func (descriptor *appDescriptorV3) GetPublisher() string {
	return descriptor.publisher
}

func (descriptor *appDescriptorV3) GetDescription() string {
	return descriptor.description
}

func (descriptor *appDescriptorV3) GetPackageVersions() map[string]*versioning.Version {
	return descriptor.packageVersions
}

func (descriptor *appDescriptorV3) GetCommandLine() []string {
	return descriptor.commandLine
}

func (descriptor *appDescriptorV3) GetIconPath() string {
	return descriptor.iconPath
}

func (descriptor *appDescriptorV3) GetSkipPackageLevels() int {
	return descriptor.skipPackageLevels
}

func (descriptor *appDescriptorV3) IsSkipUpdateCheck() bool {
	return descriptor.skipUpdateCheck
}

func (descriptor *appDescriptorV3) GetTitle() string {
	return fmt.Sprintf("%v %v", descriptor.GetName(), descriptor.GetAppVersion())
}

func (descriptor *appDescriptorV3) Init() (err error) {
	descriptor.descriptorVersion, err = versioning.ParseVersion(descriptor.DescriptorVersion)
	if err != nil {
		return fmt.Errorf("Error while parsing the Descriptor Version: %v", err.Error())
	}

	descriptor.declaredBaseURL, err = url.Parse(ensureTrailingSlash(descriptor.BaseURL))
	if err != nil {
		return fmt.Errorf("Error while parsing the Base URL: %v", err.Error())
	}

	if descriptor.DescriptorFileName != "" {
		descriptor.descriptorFileName = descriptor.DescriptorFileName
	} else {
		descriptor.descriptorFileName = defaultDescriptorFileName
	}

	descriptor.name = descriptor.Name

	descriptor.appVersion, err = versioning.ParseVersion(descriptor.Version)
	if err != nil {
		return fmt.Errorf("Error while parsing the app version: %v", err.Error())
	}

	descriptor.publisher = descriptor.Publisher
	descriptor.description = descriptor.Description

	descriptor.skipPackageLevels = descriptor.SkipPackageLevels
	descriptor.skipUpdateCheck = descriptor.SkipUpdateCheck

	if descriptor.SupportedOS != nil {
		descriptor.supportedSystems = descriptor.SupportedOS
	} else {
		descriptor.supportedSystems = []string{}
	}

	osSettings, osSettingsFound := descriptor.OS[runtime.GOOS]

	if osSettingsFound && osSettings.Packages != nil {
		descriptor.packageVersions, err = parsePackageVersions(osSettings.Packages)
	} else {
		descriptor.packageVersions, err = parsePackageVersions(descriptor.Packages)
	}

	if err != nil {
		return fmt.Errorf("Error while parsing the package versions: %v", err.Error())
	}

	if osSettingsFound && osSettings.CommandLine != nil {
		descriptor.commandLine = osSettings.CommandLine
	} else {
		descriptor.commandLine = descriptor.CommandLine
	}

	if osSettingsFound && osSettings.IconPath != "" {
		descriptor.iconPath = osSettings.IconPath
	} else {
		descriptor.iconPath = descriptor.IconPath
	}

	descriptor.actualBaseURL = getActualBaseURL(descriptor)

	return nil
}

func (descriptor *appDescriptorV3) CheckRequirements() (err error) {
	if len(descriptor.supportedSystems) > 0 {
		foundOS := false

		for _, supportedSystem := range descriptor.supportedSystems {
			if supportedSystem == runtime.GOOS {
				foundOS = true
				break
			}
		}

		if !foundOS {
			return fmt.Errorf("The current OS (%v) is not supported by %v.", runtime.GOOS, descriptor.GetTitle())
		}
	}
	return nil
}

func (descriptor *appDescriptorV3) GetFileURL(relativePath string) (fileURL *url.URL, err error) {
	return getRelativeFileURL(descriptor, relativePath)
}

func (descriptor *appDescriptorV3) GetBytes() (bytes []byte, err error) {
	return json.Marshal(*descriptor)
}
