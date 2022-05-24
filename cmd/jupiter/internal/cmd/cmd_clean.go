// Copyright 2022 Douyu
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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// Clean 清除所有的缓存
func Clean(c *cli.Context) error {
	// 1. 清除已经存在的临时模板文件
	if err := cleanTempLayout(); err != nil {
		return err
	}

	// 2. clean other ...

	color.Green("clear complete ...")
	return nil
}

// 清除已经存在的临时模板文件
func cleanTempLayout() error {
	fmt.Println("clear temp project layout ...")
	// 查看临时文件之中是否已经存在该文件夹
	tempPath := filepath.Join(os.TempDir(), "local_temp_jupiter_layout")
	// 需要刷新，提前清理缓存的文件
	if err := os.RemoveAll(tempPath); err != nil {
		return err
	}

	return nil
}
