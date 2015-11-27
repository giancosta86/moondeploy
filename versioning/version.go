/*§
  ===========================================================================
  MoonDeploy
  ===========================================================================
  Copyright (C) 2015 Gianluca Costa
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

package versioning

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major   int
	Minor   int
	Build   int
	Release int
}

func (version *Version) String() string {
	result := fmt.Sprintf("%v.%v", version.Major, version.Minor)

	if version.Build == 0 && version.Release == 0 {
		return result
	}

	result = result + "." + strconv.Itoa(version.Build)

	if version.Release == 0 {
		return result
	}

	return result + "." + strconv.Itoa(version.Release)
}

func (version *Version) CompareTo(otherVersion *Version) (result int) {
	result = version.Major - otherVersion.Major
	if result != 0 {
		return result
	}

	result = version.Minor - otherVersion.Minor
	if result != 0 {
		return result
	}

	result = version.Build - otherVersion.Build
	if result != 0 {
		return result
	}

	return version.Release - otherVersion.Release
}

func (version *Version) NewerThan(otherVersion *Version) bool {
	return version.CompareTo(otherVersion) > 0
}

func ParseVersion(versionString string) (version *Version, err error) {
	version = &Version{}

	versionComponents := strings.Split(versionString, ".")

	version.Major, err = strconv.Atoi(versionComponents[0])
	if err != nil {
		return nil, err
	}
	if version.Major < 0 {
		return nil, fmt.Errorf("Version major cannot be < 0: %v", version.Major)
	}

	if len(versionComponents) > 1 {
		version.Minor, err = strconv.Atoi(versionComponents[1])
		if err != nil {
			return nil, err
		}
		if version.Minor < 0 {
			return nil, fmt.Errorf("Version minor cannot be < 0: %v", version.Minor)
		}

		if len(versionComponents) > 2 {
			version.Build, err = strconv.Atoi(versionComponents[2])
			if err != nil {
				return nil, err
			}
			if version.Build < 0 {
				return nil, fmt.Errorf("Version build cannot be < 0: %v", version.Build)
			}

			if len(versionComponents) > 3 {
				version.Release, err = strconv.Atoi(versionComponents[3])
				if err != nil {
					return nil, err
				}
				if version.Release < 0 {
					return nil, fmt.Errorf("Version release cannot be < 0: %v", version.Release)
				}
			}
		}
	}

	return version, nil
}

func MustParseVersion(versionString string) (version *Version) {
	version, err := ParseVersion(versionString)

	if err != nil {
		panic(err)
	}

	return version
}
