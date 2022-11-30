package cmd

import (
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "init jupiter dependencies",
		Action:  Init,
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
		},
		Action: New,
	},
	{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "auto restart program when files changed",
		Action:  Run,
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
		Action:  Update,
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
		Action: Clean,
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
		Action:  Struct2Interface,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,dir",
				Usage: "please specify the code path",
				Value: ".",
			},
		},
	},
}

func run(args []string) {
	app := cli.NewApp()
	app.Usage = "Fast bootstrap tool for jupiter framework"
	app.Commands = Commands

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCMD(t *testing.T) {
	run([]string{"jupiter", "update"})
	run([]string{"jupiter", "init"})
	run([]string{"jupiter", "new", "/tmp/test-go", "-remote", "github.com/douyu/jupiter-layout", "-branch", "main"})
	run([]string{"jupiter", "struct2interface", "-d", "/tmp/test-go/internal/pkg"})
	assert.DirExists(t, "/tmp/github.com_douyu_jupiter_layout/")
	run([]string{"jupiter", "clean"})
	assert.NoDirExists(t, "/tmp/github.com_douyu_jupiter_layout/")
}
