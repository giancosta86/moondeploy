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
)

func NewAppDescriptorFromPath(descriptorPath string) (descriptor AppDescriptor, err error) {
	descriptorBytes, err := ioutil.ReadFile(descriptorPath)
	if err != nil {
		return nil, err
	}

	return NewAppDescriptorFromBytes(descriptorBytes)
}

func NewAppDescriptorFromBytes(descriptorBytes []byte) (descriptor AppDescriptor, err error) {
	descriptor, err = createV3Descriptor(descriptorBytes)

	descriptorVersion, err := descriptor.GetDescriptorVersion()
	if err != nil {
		return nil, err
	}

	switch descriptorVersion.Major {
	case 3:
		break

	case 2:
	case 1:
		descriptor, err = createV1V2Descriptor(descriptorBytes)

	default:
		return nil, fmt.Errorf("Unsupported descriptor version (%v). Please, consider updating MoonDeploy.",
			descriptorVersion)
	}

	err = descriptor.Validate()
	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

func createV3Descriptor(descriptorBytes []byte) (descriptor AppDescriptor, err error) {
	descriptor = &AppDescriptorV1V2{} //TODO: Instantiate a V3 descriptor!!!

	err = json.Unmarshal(descriptorBytes, descriptor)
	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

func createV1V2Descriptor(descriptorBytes []byte) (descriptor AppDescriptor, err error) {
	descriptor = &AppDescriptorV1V2{}

	err = json.Unmarshal(descriptorBytes, descriptor)
	if err != nil {
		return nil, err
	}

	return descriptor, nil
}
