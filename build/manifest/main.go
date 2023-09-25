/*
 *    Copyright 2023 Girish M
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) <= 1 {
		panic("no cmd specified")
	}

	manifest, err := findManifest()
	if err != nil {
		panic("failed to find manifest: " + err.Error())
	}

	cmd := os.Args[1]
	switch cmd {
	case "id":
		dumpPluginID(manifest)

	case "version":
		dumpPluginVersion(manifest)

	case "has_server":
		if manifest.HasServer() {
			fmt.Printf("true")
		}

	case "has_webapp":
		if manifest.HasWebapp() {
			fmt.Printf("true")
		}

	default:
		panic("unrecognized command: " + cmd)
	}
}

func findManifest() (*model.Manifest, error) {
	_, manifestFilePath, err := model.FindManifest(".")
	if err != nil {
		return nil, errors.Wrap(err, "failed to find manifest in current working directory")
	}
	manifestFile, err := os.Open(manifestFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %s", manifestFilePath)
	}
	defer func(manifestFile *os.File) {
		_ = manifestFile.Close()
	}(manifestFile)

	// Re-decode the manifest, disallowing unknown fields. When we write the manifest back out,
	// we don't want to accidentally clobber anything we won't preserve.
	var manifest model.Manifest
	decoder := json.NewDecoder(manifestFile)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&manifest); err != nil {
		return nil, errors.Wrap(err, "failed to parse manifest")
	}

	return &manifest, nil
}

// dumpPluginId writes the plugin id from the given manifest to standard out
func dumpPluginID(manifest *model.Manifest) {
	fmt.Printf("%s", manifest.Id)
}

// dumpPluginVersion writes the plugin version from the given manifest to standard out
func dumpPluginVersion(manifest *model.Manifest) {
	fmt.Printf("%s", manifest.Version)
}
