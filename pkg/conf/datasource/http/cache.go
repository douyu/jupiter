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

package http

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/douyu/jupiter/pkg/util/xfile"
	"github.com/douyu/jupiter/pkg/xlog"
)

// GetFileName ...
func GetFileName(cacheKey string, cacheDir string) string {
	return cacheDir + string(os.PathSeparator) + cacheKey
}

// WriteConfigToFile ...
func WriteConfigToFile(cacheKey string, cacheDir string, content string) {
	if err := xfile.MkdirIfNecessary(cacheDir); err != nil {
		xlog.Errorf("[ERROR]:faild to MkdirIfNecessary config ,value:%s ,err:%s \n", string(content), err.Error())
		return
	}
	fileName := GetFileName(cacheKey, cacheDir)
	err := ioutil.WriteFile(fileName, []byte(content), 0666)
	if err != nil {
		xlog.Errorf("[ERROR]:faild to write config  cache:%s ,value:%s ,err:%s \n", fileName, string(content), err.Error())
	}
}

// ReadConfigFromFile ...
func ReadConfigFromFile(cacheKey string, cacheDir string) (string, error) {
	fileName := GetFileName(cacheKey, cacheDir)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read config cache file:%s,err:%s! ", fileName, err.Error())
	}
	return string(b), nil
}
