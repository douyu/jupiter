package main

import (
	"log"
	"os"
	"sort"

	"github.com/douyu/jupiter/cmd/jupiter/internal/cmd"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Usage = "Fast bootstrap tool for jupiter framework"
	app.Commands = cmd.Commands

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
