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
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"

	"github.com/giancosta86/moondeploy/v3/versioning"
)

const AnyOS = "*"

type AppDescriptorV1V2 struct {
	DescriptorVersion string
	BaseURL           string

	Name      string
	Version   string
	Publisher string

	Description string
	IconPath    map[string]string

	SkipUpdateCheck   bool
	SkipPackageLevels int

	CommandLine map[string][]string

	PackageVersions map[string]string

	//
	//Computed fields
	//
	descriptorVersion *versioning.Version
	appVersion        *versioning.Version
	declaredBaseURL   *url.URL
	actualBaseURL     *url.URL
	iconPath          string
	commandLine       []string

	packageVersions map[string]*versioning.Version
}

func (descriptor *AppDescriptorV1V2) GetDescriptorVersion() *versioning.Version {
	return descriptor.descriptorVersion
}

func (descriptor *AppDescriptorV1V2) GetActualBaseURL() *url.URL {
	return descriptor.actualBaseURL
}

func (descriptor *AppDescriptorV1V2) GetDeclaredBaseURL() *url.URL {
	return descriptor.declaredBaseURL
}

func (descriptor *AppDescriptorV1V2) GetDescriptorFileName() string {
	return DefaultDescriptorFileName
}

func (descriptor *AppDescriptorV1V2) GetName() string {
	return descriptor.Name
}

func (descriptor *AppDescriptorV1V2) GetAppVersion() *versioning.Version {
	return descriptor.appVersion
}

func (descriptor *AppDescriptorV1V2) GetPublisher() string {
	return descriptor.Publisher
}

func (descriptor *AppDescriptorV1V2) GetDescription() string {
	return descriptor.Description
}

func (descriptor *AppDescriptorV1V2) GetPackageVersions() map[string]*versioning.Version {
	return descriptor.packageVersions
}

func (descriptor *AppDescriptorV1V2) GetCommandLine() []string {
	return descriptor.commandLine
}

func (descriptor *AppDescriptorV1V2) GetSkipPackageLevels() int {
	return descriptor.SkipPackageLevels
}

func (descriptor *AppDescriptorV1V2) IsSkipUpdateCheck() bool {
	return descriptor.SkipUpdateCheck
}

func (descriptor *AppDescriptorV1V2) GetIconPath() string {
	return descriptor.iconPath
}

func (descriptor *AppDescriptorV1V2) GetTitle() string {
	return fmt.Sprintf("%v %v", descriptor.Name, descriptor.Version)
}

func (descriptor *AppDescriptorV1V2) Init() (err error) {
	descriptor.descriptorVersion, err = versioning.ParseVersion(descriptor.DescriptorVersion)
	if err != nil {
		return err
	}

	descriptor.appVersion, err = versioning.ParseVersion(descriptor.Version)
	if err != nil {
		return err
	}

	err = descriptor.setDeclaredBaseURL()
	if err != nil {
		return err
	}

	descriptor.actualBaseURL = getActualBaseURL(descriptor)

	descriptor.setIconPath()

	descriptor.setCommandLine()

	descriptor.packageVersions, err = parsePackageVersions(descriptor.PackageVersions)
	if err != nil {
		return err
	}

	return nil
}

func (descriptor *AppDescriptorV1V2) setDeclaredBaseURL() (err error) {
	descriptor.declaredBaseURL, err = url.Parse(ensureTrailingSlash(descriptor.BaseURL))
	return err
}

func (descriptor *AppDescriptorV1V2) setIconPath() {
	if descriptor.IconPath == nil {
		return
	}

	osSpecificIconPath := descriptor.IconPath[runtime.GOOS]
	if osSpecificIconPath != "" {
		descriptor.iconPath = osSpecificIconPath
		return
	}

	genericIconPath := descriptor.IconPath[AnyOS]
	if genericIconPath != "" {
		descriptor.iconPath = genericIconPath
	}
}

func (descriptor *AppDescriptorV1V2) setCommandLine() {
	if descriptor.CommandLine == nil {
		return
	}

	osSpecificCommandLine := descriptor.CommandLine[runtime.GOOS]
	if osSpecificCommandLine != nil {
		descriptor.commandLine = osSpecificCommandLine
		return
	}

	genericCommandLine := descriptor.CommandLine[AnyOS]
	if genericCommandLine != nil {
		descriptor.commandLine = genericCommandLine
	}
}

func (descriptor *AppDescriptorV1V2) CheckRequirements() (err error) {
	return nil
}

func (descriptor *AppDescriptorV1V2) GetFileURL(relativePath string) (fileURL *url.URL, err error) {
	return getRelativeFileURL(descriptor, relativePath)
}

func (descriptor *AppDescriptorV1V2) GetBytes() (bytes []byte, err error) {
	return json.Marshal(*descriptor)
}
