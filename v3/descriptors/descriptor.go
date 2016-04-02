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
	"net/url"

	"github.com/giancosta86/moondeploy/v3/versioning"
)

const defaultDescriptorFileName = "App.moondeploy"

type AppDescriptor interface {
	GetDescriptorVersion() *versioning.Version
	GetDeclaredBaseURL() *url.URL
	GetActualBaseURL() *url.URL
	GetDescriptorFileName() string

	GetName() string
	GetAppVersion() *versioning.Version
	GetPublisher() string
	GetDescription() string

	GetPackageVersions() map[string]*versioning.Version
	GetCommandLine() []string
	GetSkipPackageLevels() int
	IsSkipUpdateCheck() bool

	GetIconPath() string

	GetTitle() string

	Init() (err error)
	CheckRequirements() (err error)

	GetFileURL(relativePath string) (fileURL *url.URL, err error)

	GetBytes() (bytes []byte, err error)
}
