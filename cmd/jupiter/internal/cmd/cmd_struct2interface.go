package cmd

import (
	"github.com/hnlq715/struct2interface"
	"github.com/urfave/cli"
)

func Struct2Interface(c *cli.Context) error {
	return struct2interface.MakeDir(c.String("dir"))
}
