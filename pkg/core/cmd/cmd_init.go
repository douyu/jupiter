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
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func Init(c *cli.Context) error {

	deps := []string{
		"github.com/google/wire/cmd/wire@v0.5.0",
		"github.com/vektra/mockery/v2@v2.14.0",
		"github.com/bufbuild/buf/cmd/buf@v1.6.0",
		"github.com/onsi/ginkgo/v2/ginkgo@v2.1.3",
		"github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.8.7",
	}

	for _, dep := range deps {
		err := goinstall(dep)
		if err != nil {
			return err
		}
	}

	color.Green("jupiter init success.")
	return nil
}

func goinstall(path string) error {
	cmd := exec.Command("go", "install", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		color.Red("install %s failed, please install it manually", path)
		return err
	}
	color.Green("install %s success.", path)
	return nil
}
