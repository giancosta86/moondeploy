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
	"path/filepath"
	"strings"

	"github.com/giancosta86/moondeploy/v3/descriptors"
)

type AppGallery struct {
	Directory string
}

func NewAppGallery(galleryPath string) (appGallery *AppGallery) {
	appGallery = &AppGallery{
		Directory: galleryPath,
	}

	return appGallery
}

func (appGallery *AppGallery) GetApp(bootDescriptor descriptors.AppDescriptor) (app *App, err error) {
	appDir, err := appGallery.resolveAppDir(bootDescriptor)
	if err != nil {
		return nil, err
	}

	return &App{
		Directory:      appDir,
		bootDescriptor: bootDescriptor,
		filesDirectory: filepath.Join(appDir, filesDirName),
	}, nil
}

func (appGallery *AppGallery) resolveAppDir(bootDescriptor descriptors.AppDescriptor) (appDir string, err error) {
	baseURL := bootDescriptor.GetDeclaredBaseURL()

	hostComponent := strings.Replace(baseURL.Host, ":", "_", -1)

	appDirComponents := []string{
		appGallery.Directory,
		hostComponent}

	trimmedBasePath := strings.Trim(baseURL.Path, "/")
	baseComponents := strings.Split(trimmedBasePath, "/")

	appDirComponents = append(appDirComponents, baseComponents...)

	if hostComponent == "github.com" &&
		len(appDirComponents) > 2 &&
		appDirComponents[len(appDirComponents)-2] == "releases" &&
		appDirComponents[len(appDirComponents)-1] == "latest" {
		appDirComponents = appDirComponents[0 : len(appDirComponents)-2]
	}

	appDir = filepath.Join(appDirComponents...)

	return appDir, nil
}
