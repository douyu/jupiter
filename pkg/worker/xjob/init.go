package job

import (
	"github.com/douyu/jupiter/pkg/flag"
)

func init() {
	flag.Register(&flag.BoolFlag{
		Name:    "disable-job",
		Usage:   "--disable-job, disable job",
		Default: false,
	})
	flag.Register(&flag.StringFlag{
		Name:  "job",
		Usage: "--job, job name",
	})
}
