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

package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/core/cmd"
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "init jupiter dependencies",
		Action:  cmd.Init,
	},
	{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "generate code framework",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "remote",
				Usage: "choose remote template",
				Value: "github.com/douyu/jupiter-layout",
			},
			cli.StringFlag{
				Name:  "branch",
				Usage: "choose branch of remote template",
				Value: "main",
			},
			cli.BoolFlag{
				Name:  "upgrade",
				Usage: "upgrade remote template",
			},
			cli.StringFlag{
				Name:  "app",
				Usage: "app name",
			},
		},
		Action: cmd.New,
	},
	{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "auto restart program when files changed",
		Action:  cmd.Run,
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
		Action:  cmd.Update,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "remote",
				Usage: "choose remote repo",
				Value: "github.com/douyu/jupiter/cmd/jupiter",
			},
		},
	},
	{
		Name:   "clean",
		Usage:  "clear all cached",
		Action: cmd.Clean,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "remote",
				Usage: "choose remote template",
				Value: "github.com/douyu/jupiter-layout",
			},
		},
	},
	{
		Name:    "struct2interface",
		Aliases: []string{"struct2interface"},
		Usage:   "Auto generate interface from struct for golang",
		Action:  cmd.Struct2Interface,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,dir",
				Usage: "please specify the code path",
				Value: ".",
			},
		},
	},
	{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "show version",
		Action: func(c *cli.Context) error {
			fmt.Printf("jupiter %v\n", pkg.JupiterVersion())
			return nil
		},
	},
}

func main() {

	app := cli.NewApp()
	app.Usage = "Fast bootstrap tool for jupiter framework"
	app.Commands = Commands

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
