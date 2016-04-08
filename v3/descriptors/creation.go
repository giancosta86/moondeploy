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
	"io/ioutil"

	"github.com/giancosta86/moondeploy"
	"github.com/giancosta86/moondeploy/v3/log"
	"github.com/giancosta86/moondeploy/v3/versioning"
)

type BasicDescriptor struct {
	DescriptorVersion string
}

func NewAppDescriptorFromPath(descriptorPath string) (descriptor AppDescriptor, err error) {
	descriptorBytes, err := ioutil.ReadFile(descriptorPath)
	if err != nil {
		return nil, err
	}

	return NewAppDescriptorFromBytes(descriptorBytes)
}

func NewAppDescriptorFromBytes(descriptorBytes []byte) (descriptor AppDescriptor, err error) {
	basicDescriptor, err := createBasicDescriptor(descriptorBytes)
	if err != nil {
		return nil, err
	}

	descriptorVersion, err := versioning.ParseVersion(basicDescriptor.DescriptorVersion)
	if err != nil {
		return nil, fmt.Errorf("Invalid Descriptor Version: %v", err.Error())
	}

	switch descriptorVersion.Major {
	case 3:
		log.Notice("V3 descriptor found! Deserializing it")
		descriptor, err = createV3Descriptor(descriptorBytes)

	case 2:
	case 1:
		log.Notice("V1/V2 descriptor found! Deserializing it")
		descriptor, err = createV1V2Descriptor(descriptorBytes)

	default:
		return nil, fmt.Errorf("Unsupported descriptor version (%v). Please, consider updating MoonDeploy - your current version is %v.",
			descriptorVersion,
			moondeploy.Version)
	}

	err = descriptor.Init()
	if err != nil {
		return nil, err
	}

	err = validate(descriptor)
	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

func createBasicDescriptor(descriptorBytes []byte) (descriptor *BasicDescriptor, err error) {
	descriptor = &BasicDescriptor{}

	err = json.Unmarshal(descriptorBytes, descriptor)
	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

func createV3Descriptor(descriptorBytes []byte) (descriptor AppDescriptor, err error) {
	descriptor = &appDescriptorV3{}

	err = json.Unmarshal(descriptorBytes, descriptor)
	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

func createV1V2Descriptor(descriptorBytes []byte) (descriptor AppDescriptor, err error) {
	descriptor = &appDescriptorV1V2{}

	err = json.Unmarshal(descriptorBytes, descriptor)
	if err != nil {
		return nil, err
	}

	return descriptor, nil
}
