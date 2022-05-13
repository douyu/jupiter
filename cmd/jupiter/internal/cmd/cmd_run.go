package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/douyu/jupiter/cmd/jupiter/internal/runner"
	"github.com/urfave/cli"
)

// Run 运行程序
func Run(c *cli.Context) error {
	debugMode := c.Bool("debug")
	cfgPath := c.String("c")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var err error
	r, err := runner.NewEngine(cfgPath, debugMode)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go func() {
		<-sigs
		r.Stop()
	}()

	r.Run()

	return nil
}
