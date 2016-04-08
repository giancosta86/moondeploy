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
	"strings"
)

func validate(descriptor AppDescriptor) (err error) {
	if descriptor.GetDescriptorVersion() == nil {
		return fmt.Errorf("Descriptor Version field is missing")
	}

	if descriptor.GetDeclaredBaseURL() == nil {
		return fmt.Errorf("Declared Base URL field is missing")
	}

	if descriptor.GetActualBaseURL() == nil {
		return fmt.Errorf("Actual Base URL field is missing")
	}

	if strings.TrimSpace(descriptor.GetDescriptorFileName()) == "" {
		return fmt.Errorf("Descriptor File Name field is missing")
	}

	if strings.TrimSpace(descriptor.GetName()) == "" {
		return fmt.Errorf("Name field is missing")
	}

	if descriptor.GetAppVersion() == nil {
		return fmt.Errorf("App version field is missing")
	}

	if strings.TrimSpace(descriptor.GetPublisher()) == "" {
		return fmt.Errorf("Publisher field is missing")
	}

	if strings.TrimSpace(descriptor.GetDescription()) == "" {
		return fmt.Errorf("Description field is missing")
	}

	if descriptor.GetPackageVersions() == nil {
		return fmt.Errorf("Package versions field is missing")
	}

	if descriptor.GetCommandLine() == nil || len(descriptor.GetCommandLine()) == 0 {
		return fmt.Errorf("Command Line field is missing")
	}

	if descriptor.GetSkipPackageLevels() < 0 {
		return fmt.Errorf("SkipPackageLevels field must be >= 0")
	}

	if strings.TrimSpace(descriptor.GetTitle()) == "" {
		return fmt.Errorf("The title is missing")
	}

	return nil
}

func CheckDescriptorMatch(descriptor AppDescriptor, otherDescriptor AppDescriptor) (err error) {
	if descriptor.GetDeclaredBaseURL().String() != otherDescriptor.GetDeclaredBaseURL().String() {
		return fmt.Errorf("The descriptors have different BaseURL's:\n\t'%v'\n\t'%v'",
			descriptor.GetDeclaredBaseURL(),
			otherDescriptor.GetDeclaredBaseURL())
	}

	if descriptor.GetDescriptorFileName() != otherDescriptor.GetDescriptorFileName() {
		return fmt.Errorf("The descriptors have different Descriptor File Name values:\n\t'%v'\n\t'%v",
			descriptor.GetDescriptorFileName(),
			otherDescriptor.GetDescriptorFileName())
	}

	if descriptor.GetName() != otherDescriptor.GetName() {
		return fmt.Errorf("The descriptors have different Name values:\n\t'%v'\n\t'%v",
			descriptor.GetName(),
			otherDescriptor.GetName())
	}

	return nil
}
