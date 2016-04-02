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

	"github.com/giancosta86/moondeploy/versioning"
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

	//Cache fields
	version         *versioning.Version
	declaredBaseURL *url.URL
	actualBaseURL   *url.URL
	iconPath        string
	commandLine     []string

	packageVersions map[string]*versioning.Version
}

func (descriptor *AppDescriptorV1V2) GetDescriptorVersion() (*versioning.Version, error) {
	return versioning.ParseVersion(descriptor.DescriptorVersion)
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
	return descriptor.version
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

func (descriptor *AppDescriptorV1V2) Validate() (err error) {
	if descriptor.BaseURL == "" {
		return fmt.Errorf("Base URL field is missing")
	}

	if descriptor.Name == "" {
		return fmt.Errorf("Name field is missing")
	}

	if descriptor.Version == "" {
		return fmt.Errorf("Version field is missing")
	}

	if descriptor.Publisher == "" {
		return fmt.Errorf("Publisher field is missing")
	}

	descriptor.version, err = versioning.ParseVersion(descriptor.Version)
	if err != nil {
		return err
	}

	err = descriptor.setDeclaredBaseURL()
	if err != nil {
		return err
	}

	descriptor.actualBaseURL = getActualBaseURL(descriptor)

	if descriptor.IconPath == nil {
		descriptor.IconPath = make(map[string]string)
	}

	descriptor.setIconPath()

	if descriptor.SkipPackageLevels < 0 {
		return fmt.Errorf("SkipPackageLevels field must be >= 0")
	}

	if descriptor.CommandLine == nil {
		descriptor.CommandLine = make(map[string][]string)
	}

	descriptor.setCommandLine()

	err = descriptor.setPackageVersions()
	if err != nil {
		return err
	}

	return nil
}

func (descriptor *AppDescriptorV1V2) setDeclaredBaseURL() (err error) {
	if descriptor.BaseURL[len(descriptor.BaseURL)-1] != '/' {
		descriptor.BaseURL = descriptor.BaseURL + "/"
	}

	descriptor.declaredBaseURL, err = url.Parse(descriptor.BaseURL)

	return err
}

func (descriptor *AppDescriptorV1V2) setIconPath() {
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
	osSpecificCommandLine := descriptor.CommandLine[runtime.GOOS]
	if osSpecificCommandLine != nil {
		descriptor.commandLine = osSpecificCommandLine
	}

	genericCommandLine := descriptor.CommandLine[AnyOS]
	if genericCommandLine != nil {
		descriptor.commandLine = genericCommandLine
	}
}

func (descriptor *AppDescriptorV1V2) setPackageVersions() (err error) {
	if descriptor.PackageVersions == nil {
		descriptor.PackageVersions = make(map[string]string)
	}

	descriptor.packageVersions = make(map[string]*versioning.Version)
	for packageName, packageVersionString := range descriptor.PackageVersions {
		if packageVersionString != "" {
			descriptor.packageVersions[packageName], err = versioning.ParseVersion(packageVersionString)
			if err != nil {
				return fmt.Errorf("Invalid version string for package '%v': '%v'",
					packageName,
					packageVersionString)
			}
		} else {
			descriptor.packageVersions[packageName] = nil
		}
	}

	return nil
}

func (descriptor *AppDescriptorV1V2) CheckMatch(otherDescriptor AppDescriptor) (err error) {
	if descriptor.GetName() != otherDescriptor.GetName() {
		return fmt.Errorf("The descriptors have different Name values:\n\t'%v'\n\t'%v",
			descriptor.GetName(),
			otherDescriptor.GetName())
	}

	if descriptor.GetDescriptorFileName() != otherDescriptor.GetDescriptorFileName() {
		return fmt.Errorf("The descriptors have different DescriptorFileName values:\n\t'%v'\n\t'%v",
			descriptor.GetDescriptorFileName(),
			otherDescriptor.GetDescriptorFileName())
	}

	if descriptor.GetDeclaredBaseURL().String() != otherDescriptor.GetDeclaredBaseURL().String() {
		return fmt.Errorf("The descriptors have different BaseURL's:\n\t'%v'\n\t'%v'",
			descriptor.GetDeclaredBaseURL(),
			otherDescriptor.GetDeclaredBaseURL())
	}

	return nil
}

func (descriptor *AppDescriptorV1V2) CheckRequirements() (err error) {
	if descriptor.commandLine == nil {
		return fmt.Errorf("The app does does not provide a command line for this operating system: %v", runtime.GOOS)
	}

	return nil
}

func (descriptor *AppDescriptorV1V2) GetFileURL(relativePath string) (fileURL *url.URL, err error) {
	return getRelativeFileURL(descriptor, relativePath)
}

func (descriptor *AppDescriptorV1V2) GetBytes() (bytes []byte, err error) {
	return json.Marshal(*descriptor)
}
