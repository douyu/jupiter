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
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

// Update 更新到最新版本
func Update(c *cli.Context) error {

	remote := c.String("remote")

	update := fmt.Sprintf("go install %s@latest\n", remote)

	if runtime.Version() < "go1.18" {
		fmt.Println("当前安装的golang版本小于1.18，请升级！")
		return nil
	}

	cmds := []string{update}

	for _, cmd := range cmds {
		fmt.Println(cmd)
		cmd := exec.Command("bash", "-c", cmd)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
