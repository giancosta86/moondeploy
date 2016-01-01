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
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/giancosta86/caravel"

	"github.com/giancosta86/moondeploy"
	"github.com/giancosta86/moondeploy/custom"
	"github.com/giancosta86/moondeploy/logging"
	"github.com/giancosta86/moondeploy/versioning"
)

const AnyOS = "*"

type AppDescriptor struct {
	DescriptorVersion *versioning.Version

	Name            string
	Version         *versioning.Version
	SkipUpdateCheck bool

	Description string
	IconPath    map[string]string

	BaseURL   *url.URL
	Publisher string

	PackageVersions map[string]*versioning.Version

	CommandLine map[string][]string

	SkipPackageLevels int
}

func NewAppDescriptorFromPath(descriptorPath string) (descriptor *AppDescriptor, err error) {
	descriptorBytes, err := ioutil.ReadFile(descriptorPath)
	if err != nil {
		return nil, err
	}

	return NewAppDescriptorFromBytes(descriptorBytes)
}

func NewAppDescriptorFromBytes(descriptorBytes []byte) (descriptor *AppDescriptor, err error) {
	rawDescriptor := &rawAppDescriptor{}

	err = json.Unmarshal(descriptorBytes, rawDescriptor)
	if err != nil {
		return nil, err
	}

	fullDescriptor, err := rawDescriptor.toFull()
	if err != nil {
		return nil, err
	}

	return fullDescriptor, nil
}

func (descriptor *AppDescriptor) GetTitle() string {
	return fmt.Sprintf("%v %v", descriptor.Name, descriptor.Version)
}

func (descriptor *AppDescriptor) GetActualIconPath(appFilesDir string) string {
	osSpecificIconPath := descriptor.IconPath[runtime.GOOS]
	if osSpecificIconPath != "" {
		return filepath.Join(appFilesDir, osSpecificIconPath)
	}

	genericIconPath := descriptor.IconPath[AnyOS]
	if genericIconPath != "" {
		return filepath.Join(appFilesDir, genericIconPath)
	}

	return moondeploy.GetIconPath()
}

func (descriptor *AppDescriptor) GetCommandLine() (commandLine []string, err error) {
	osSpecificCommandLine := descriptor.CommandLine[runtime.GOOS]
	if osSpecificCommandLine != nil {
		return osSpecificCommandLine, nil
	}

	genericCommandLine := descriptor.CommandLine[AnyOS]
	if genericCommandLine != nil {
		return genericCommandLine, nil
	}

	return nil, fmt.Errorf("The app does does not support this operating system: %v", runtime.GOOS)
}

func (descriptor *AppDescriptor) Validate() (err error) {
	if descriptor.DescriptorVersion == nil {
		return fmt.Errorf("Descriptor version is missing")
	}

	if descriptor.DescriptorVersion.NewerThan(moondeploy.Version) {
		return fmt.Errorf("Descriptor version (%v) is newer than the current MoonDeploy version (%v). Please, consider updating MoonDeploy.", descriptor.DescriptorVersion, moondeploy.Version)
	}

	if descriptor.Name == "" {
		return fmt.Errorf("Name is missing")
	}

	if descriptor.Version == nil {
		return fmt.Errorf("Version is missing")
	}

	if descriptor.BaseURL == nil {
		return fmt.Errorf("Base URL is missing")
	}

	if descriptor.Publisher == "" {
		return fmt.Errorf("Publisher is missing")
	}

	if descriptor.SkipPackageLevels < 0 {
		return fmt.Errorf("SkipPackageLevels must be >= 0")
	}

	return nil
}

func (descriptor *AppDescriptor) GetBaseFileURL(relativePath string) (fileURL *url.URL, err error) {
	if path.IsAbs(relativePath) {
		return nil, fmt.Errorf("Absolute paths are not allowed: '%v'", relativePath)
	}

	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return nil, err
	}

	return descriptor.BaseURL.ResolveReference(relativeURL), nil
}

func (descriptor *AppDescriptor) ToBytes() (bytes []byte, err error) {
	rawDescriptor := descriptor.toRaw()

	return json.Marshal(*rawDescriptor)
}

func (descriptor *AppDescriptor) toRaw() (rawDescriptor *rawAppDescriptor) {
	rawDescriptor = &rawAppDescriptor{}

	rawDescriptor.DescriptorVersion = descriptor.DescriptorVersion.String()

	rawDescriptor.Name = descriptor.Name
	rawDescriptor.Version = descriptor.Version.String()
	rawDescriptor.SkipUpdateCheck = descriptor.SkipUpdateCheck

	rawDescriptor.Description = descriptor.Description
	rawDescriptor.IconPath = descriptor.IconPath

	rawDescriptor.BaseURL = descriptor.BaseURL.String()
	rawDescriptor.Publisher = descriptor.Publisher

	rawDescriptor.PackageVersions = make(map[string]string)
	for packageName, packageVersion := range descriptor.PackageVersions {
		var packageVersionString string
		if packageVersion != nil {
			packageVersionString = packageVersion.String()
		} else {
			packageVersionString = ""
		}

		rawDescriptor.PackageVersions[packageName] = packageVersionString
	}

	rawDescriptor.CommandLine = descriptor.CommandLine

	rawDescriptor.SkipPackageLevels = descriptor.SkipPackageLevels

	return rawDescriptor
}

func (descriptor *AppDescriptor) CheckMatch(otherDescriptor *AppDescriptor) (err error) {
	if descriptor.Name != otherDescriptor.Name {
		return fmt.Errorf("The descriptors have different Name:\n\t'%v'\n\t'%v", descriptor.Name, otherDescriptor.Name)
	}

	if descriptor.BaseURL.String() != otherDescriptor.BaseURL.String() {
		return fmt.Errorf("The descriptors have different BaseURL:\n\t'%v'\n\t'%v'", descriptor.BaseURL, otherDescriptor.BaseURL)
	}

	return nil
}

func (remoteDescriptor *AppDescriptor) GetPackagesToUpdate(localDescriptor *AppDescriptor) []string {
	if localDescriptor == nil {
		packagesToUpdate := []string{}

		for packageName, _ := range remoteDescriptor.PackageVersions {
			packagesToUpdate = append(packagesToUpdate, packageName)
		}

		return packagesToUpdate
	}

	if !remoteDescriptor.Version.NewerThan(localDescriptor.Version) {
		return []string{}
	}

	packagesToUpdate := []string{}

	for remotePackageName, remotePackageVersion := range remoteDescriptor.PackageVersions {
		localPackageVersion := localDescriptor.PackageVersions[remotePackageName]

		if remotePackageVersion == nil ||
			localPackageVersion == nil ||
			remotePackageVersion.NewerThan(localPackageVersion) {
			packagesToUpdate = append(packagesToUpdate, remotePackageName)
		}
	}

	return packagesToUpdate
}

func (descriptor *AppDescriptor) InstallPackage(
	packageName string,
	targetFilesDir string,
	settings *custom.Settings,
	progressCallback caravel.RetrievalProgressCallback) (err error) {

	packageURL, err := descriptor.GetBaseFileURL(packageName)
	if err != nil {
		return err
	}

	logging.Info("Creating package temp file...")
	packageTempFile, err := ioutil.TempFile(os.TempDir(), packageName)
	if err != nil {
		return err
	}
	packageTempFilePath := packageTempFile.Name()
	logging.Info("Package temp file created '%v'", packageTempFilePath)

	defer func() {
		packageTempFile.Close()

		logging.Info("Deleting package temp file: '%v'", packageTempFilePath)
		tempFileRemovalErr := os.Remove(packageTempFilePath)
		if tempFileRemovalErr != nil {
			logging.Warning("Could not remove the package temp file! '%v'", tempFileRemovalErr)
		} else {
			logging.Notice("Package temp file removed")
		}
	}()

	logging.Info("Retrieving package: %v", packageURL)
	err = caravel.RetrieveChunksFromURL(packageURL, packageTempFile, settings.BufferSize, progressCallback)
	if err != nil {
		return err
	}
	logging.Notice("Package retrieved")

	logging.Info("Closing the package temp file...")
	packageTempFile.Close()
	if err != nil {
		return err
	}
	logging.Notice("Package temp file closed")

	err = os.MkdirAll(targetFilesDir, 0700)
	if err != nil {
		return err
	}

	logging.Info("Extracting the package. Skipping levels: %v...", descriptor.SkipPackageLevels)
	err = caravel.ExtractZipSkipLevels(packageTempFilePath, targetFilesDir, descriptor.SkipPackageLevels)
	if err != nil {
		return err
	}
	logging.Notice("Package extracted")

	return nil
}
