package cmd

import (
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func Init(c *cli.Context) error {

	err := goinstall("github.com/google/wire/cmd/wire@v0.5.0")
	if err != nil {
		return err
	}

	err = goinstall("github.com/vektra/mockery/v2@v2.14.0")
	if err != nil {
		return err
	}

	err = goinstall("github.com/bufbuild/buf/cmd/buf@v1.6.0")
	if err != nil {
		return err
	}

	err = goinstall("github.com/onsi/ginkgo/v2/ginkgo@v2.1.3")
	if err != nil {
		return err
	}

	err = goinstall("github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.8.7")
	if err != nil {
		return err
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
