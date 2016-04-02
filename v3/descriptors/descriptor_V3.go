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

type AppDescriptorV3 struct {
	DescriptorVersion string

	OsSettings

	OS map[string]OsSettings

	//
	//Computed fields
	//
	descriptorVersion  *versioning.Version
	declaredBaseURL    *url.URL
	actualBaseURL      *url.URL
	descriptorFileName string

	name        string
	appVersion  *versioning.Version
	publisher   string
	description string

	packageVersions map[string]*versioning.Version
	commandLine     []string

	iconPath string

	skipPackageLevels int
	skipUpdateCheck   bool
}

type OsSettings struct {
	BaseURL            string
	DescriptorFileName string

	Name        string
	Version     string
	Publisher   string
	Description string

	PackageVersions map[string]string
	CommandLine     []string
	IconPath        string

	SkipPackageLevels int
	SkipUpdateCheck   bool
}

func (descriptor *AppDescriptorV3) GetDescriptorVersion() *versioning.Version {
	return descriptor.descriptorVersion
}

func (descriptor *AppDescriptorV3) GetDeclaredBaseURL() *url.URL {
	return descriptor.declaredBaseURL
}

func (descriptor *AppDescriptorV3) GetActualBaseURL() *url.URL {
	return descriptor.actualBaseURL
}

func (descriptor *AppDescriptorV3) GetDescriptorFileName() string {
	return descriptor.descriptorFileName
}

func (descriptor *AppDescriptorV3) GetName() string {
	return descriptor.name
}

func (descriptor *AppDescriptorV3) GetAppVersion() *versioning.Version {
	return descriptor.appVersion
}

func (descriptor *AppDescriptorV3) GetPublisher() string {
	return descriptor.publisher
}

func (descriptor *AppDescriptorV3) GetDescription() string {
	return descriptor.description
}

func (descriptor *AppDescriptorV3) GetPackageVersions() map[string]*versioning.Version {
	return descriptor.packageVersions
}

func (descriptor *AppDescriptorV3) GetCommandLine() []string {
	return descriptor.commandLine
}

func (descriptor *AppDescriptorV3) GetIconPath() string {
	return descriptor.iconPath
}

func (descriptor *AppDescriptorV3) GetSkipPackageLevels() int {
	return descriptor.skipPackageLevels
}

func (descriptor *AppDescriptorV3) IsSkipUpdateCheck() bool {
	return descriptor.skipUpdateCheck
}

func (descriptor *AppDescriptorV3) GetTitle() string {
	return fmt.Sprintf("%v %v", descriptor.GetName(), descriptor.GetAppVersion())
}

func (descriptor *AppDescriptorV3) Init() (err error) {
	descriptor.descriptorVersion, err = versioning.ParseVersion(descriptor.DescriptorVersion)
	if err != nil {
		return err
	}

	osSettings, osSettingsFound := descriptor.OS[runtime.GOOS]

	if osSettingsFound && osSettings.BaseURL != "" {
		descriptor.declaredBaseURL, err = url.Parse(ensureTrailingSlash(osSettings.BaseURL))
	} else {
		descriptor.declaredBaseURL, err = url.Parse(ensureTrailingSlash(descriptor.BaseURL))
	}

	if err != nil {
		return
	}

	if osSettingsFound && osSettings.DescriptorFileName != "" {
		descriptor.descriptorFileName = osSettings.DescriptorFileName
	} else {
		descriptor.descriptorFileName = descriptor.DescriptorFileName
	}

	if osSettingsFound && descriptor.descriptorFileName == "" {
		descriptor.descriptorFileName = defaultDescriptorFileName
	}

	if osSettingsFound && osSettings.Name != "" {
		descriptor.name = osSettings.Name
	} else {
		descriptor.name = descriptor.Name
	}

	if osSettingsFound && osSettings.Version != "" {
		descriptor.appVersion, err = versioning.ParseVersion(osSettings.Version)
	} else {
		descriptor.appVersion, err = versioning.ParseVersion(descriptor.Version)
	}

	if err != nil {
		return err
	}

	if osSettingsFound && osSettings.Publisher != "" {
		descriptor.publisher = osSettings.Publisher
	} else {
		descriptor.publisher = descriptor.Publisher
	}

	if osSettingsFound && osSettings.Description != "" {
		descriptor.description = osSettings.Description
	} else {
		descriptor.description = descriptor.Description
	}

	if osSettingsFound && osSettings.PackageVersions != nil {
		descriptor.packageVersions, err = parsePackageVersions(osSettings.PackageVersions)
	} else {
		descriptor.packageVersions, err = parsePackageVersions(descriptor.PackageVersions)
	}

	if err != nil {
		return err
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

	if osSettingsFound {
		descriptor.skipPackageLevels = osSettings.SkipPackageLevels
	} else {
		descriptor.skipPackageLevels = osSettings.SkipPackageLevels
	}

	if osSettingsFound {
		descriptor.skipUpdateCheck = osSettings.SkipUpdateCheck
	} else {
		descriptor.skipUpdateCheck = osSettings.SkipUpdateCheck
	}

	descriptor.actualBaseURL = getActualBaseURL(descriptor)

	return nil
}

func (descriptor *AppDescriptorV3) CheckRequirements() (err error) {
	return nil
}

func (descriptor *AppDescriptorV3) GetFileURL(relativePath string) (fileURL *url.URL, err error) {
	return getRelativeFileURL(descriptor, relativePath)
}

func (descriptor *AppDescriptorV3) GetBytes() (bytes []byte, err error) {
	return json.Marshal(*descriptor)
}
