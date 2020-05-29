// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xfile

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/multierr"
)

// GetCurrentDirectory ...
func GetCurrentDirectory() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("err", err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// GetCurrentPackage ...
func GetCurrentPackage() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("err", err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// MakeDirectory ...
func MakeDirectory(dirs ...string) error {
	var errs error
	for _, dir := range dirs {
		if !Exists(dir) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				errs = multierr.Append(errs, err)
			}
		}
	}
	return errs
}
