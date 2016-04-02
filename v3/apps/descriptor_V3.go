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

type OsSettingsV3 struct {
	RemoteURL         string
	PackageVersions   map[string]string
	CommandLine       string
	SkipPackageLevels int
	SkipUpdateCheck   bool
	Description       string
	IconPath          string
}

type AppDescriptorV3 struct {
	DescriptorVersion string

	Name      string
	Version   string
	Publisher string

	OsSettingsV3

	Systems map[string]OsSettingsV3
}