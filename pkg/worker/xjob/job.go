package job

import (
	"github.com/douyu/jupiter/pkg/flag"
)

func init() {
	flag.Register(
		&flag.StringFlag{
			Name:    "job",
			Usage:   "--job",
			Default: "",
		},
	)
}

// Runner ...
type Runner interface {
	Run()
}
