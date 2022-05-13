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
