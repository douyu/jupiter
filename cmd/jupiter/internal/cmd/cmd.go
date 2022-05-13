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

import "github.com/urfave/cli"

var Commands = []cli.Command{
	{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "generate code framework",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "type",
				Usage: "choose template",
				Value: "jupiter-layout",
			},
		},
		Action: func(c *cli.Context) error {
			return New(c)
		},
	},
	{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "auto restart program when files changed",
		Action: func(c *cli.Context) error {
			return Run(c)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "c",
				Value: ".jupiter.toml",
				Usage: "指定启动配置文件",
			},
			cli.BoolFlag{
				Name:  "debug",
				Usage: "debug mode",
			},
		},
	},
	{
		Name:    "update",
		Aliases: []string{"upgrade"},
		Usage:   "Upgrade to the latest version",
		Action: func(c *cli.Context) error {
			return Update(c)
		},
	},
}
