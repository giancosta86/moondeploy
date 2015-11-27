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

package test

import (
	"testing"

	"github.com/giancosta86/moondeploy/apps"
	"github.com/giancosta86/moondeploy/versioning"
)

func testGetPackagesToUpdate(
	t *testing.T,

	remoteVersion string,
	localVersion string,

	remotePackages map[string]string,
	localPackages map[string]string,

	expectedPackagesToUpdate []string) {

	actualRemotePackages := make(map[string]*versioning.Version)
	for packageName, versionString := range remotePackages {
		if versionString != "" {
			actualRemotePackages[packageName] = versioning.MustParseVersion(versionString)
		} else {
			actualRemotePackages[packageName] = nil
		}
	}

	actualLocalPackages := make(map[string]*versioning.Version)
	for packageName, versionString := range localPackages {
		if versionString != "" {
			actualLocalPackages[packageName] = versioning.MustParseVersion(versionString)
		} else {
			actualLocalPackages[packageName] = nil
		}
	}

	remoteDescriptor := &apps.AppDescriptor{
		Version:         versioning.MustParseVersion(remoteVersion),
		PackageVersions: actualRemotePackages,
	}

	localDescriptor := &apps.AppDescriptor{
		Version:         versioning.MustParseVersion(localVersion),
		PackageVersions: actualLocalPackages,
	}

	packagesToUpdate := remoteDescriptor.GetPackagesToUpdate(localDescriptor)

	if len(packagesToUpdate) != len(expectedPackagesToUpdate) {
		t.Fatalf("Expected packages to update: %v. Found: %v", len(expectedPackagesToUpdate), len(packagesToUpdate))
	}

	packagesToUpdateSet := make(map[string]bool)
	for _, packageToUpdate := range packagesToUpdate {
		packagesToUpdateSet[packageToUpdate] = true
	}

	for _, expectedPackage := range expectedPackagesToUpdate {
		_, expectedPackageFound := packagesToUpdateSet[expectedPackage]

		if !expectedPackageFound {
			t.Fatalf("Expected package to update '%v' not found", expectedPackage)
		}
	}
}

func TestUpdateWithDifferentDescriptorVersions_1(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "",
		},

		map[string]string{
			"alpha": "",
		},

		[]string{"alpha"})
}

func TestUpdateWithDifferentDescriptorVersions_2(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "1.0",
		},

		map[string]string{
			"alpha": "",
		},

		[]string{"alpha"})
}

func TestUpdateWithDifferentDescriptorVersions_3(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "",
		},

		map[string]string{
			"alpha": "2.0",
		},

		[]string{"alpha"})
}

func TestUpdateWithDifferentDescriptorVersions_4(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "1.1",
		},

		map[string]string{
			"alpha": "1.0",
		},

		[]string{"alpha"})
}

func TestUpdateWithDifferentDescriptorVersions_5(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "1.4",
		},

		map[string]string{
			"alpha": "1.4",
		},

		[]string{})
}

func TestUpdateWithDifferentDescriptorVersions_6(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "",
			"beta":  "1.3",
		},

		map[string]string{
			"alpha": "",
		},

		[]string{"alpha", "beta"})
}

func TestUpdateWithDifferentDescriptorVersions_7(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "1.0",
			"beta":  "1.4",
		},

		map[string]string{
			"alpha": "",
			"beta":  "1.3",
		},

		[]string{"alpha", "beta"})
}

func TestUpdateWithDifferentDescriptorVersions_8(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "",
			"beta":  "",
		},

		map[string]string{
			"alpha": "2.0",
			"beta":  "1.2",
		},

		[]string{"alpha", "beta"})
}

func TestUpdateWithDifferentDescriptorVersions_9(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "1.1",
			"beta":  "1.2",
		},

		map[string]string{
			"alpha": "1.0",
			"beta":  "",
		},

		[]string{"alpha", "beta"})
}

func TestUpdateWithDifferentDescriptorVersions_10(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.1",
		"1.0",

		map[string]string{
			"alpha": "1.4",
		},

		map[string]string{
			"alpha": "1.4",
			"beta":  "1.2",
		},

		[]string{})
}

func TestUpdateWithEqualDescriptorVersions_1(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.0",
		"1.0",

		map[string]string{
			"alpha": "",
		},

		map[string]string{
			"alpha": "",
		},

		[]string{})
}

func TestUpdateWithEqualDescriptorVersions_2(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.0",
		"1.0",

		map[string]string{
			"alpha": "1.0",
		},

		map[string]string{
			"alpha": "",
		},

		[]string{})
}

func TestUpdateWithEqualDescriptorVersions_3(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.0",
		"1.0",

		map[string]string{
			"alpha": "",
		},

		map[string]string{
			"alpha": "2.0",
		},

		[]string{})
}

func TestUpdateWithEqualDescriptorVersions_4(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.0",
		"1.0",

		map[string]string{
			"alpha": "1.1",
		},

		map[string]string{
			"alpha": "1.0",
		},

		[]string{})
}

func TestUpdateWithEqualDescriptorVersions_5(t *testing.T) {
	testGetPackagesToUpdate(
		t,

		"1.0",
		"1.0",

		map[string]string{
			"alpha": "1.4",
		},

		map[string]string{
			"alpha": "1.4",
		},

		[]string{})
}
