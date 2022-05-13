package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

// Update 更新到最新版本
func Update(c *cli.Context) error {
	update := "go install github.com/douyu/jupiter/cmd/jupiter@latest\n"

	if runtime.Version() < "go1.16" {
		fmt.Println("当前安装的golang版本小于1.16，请升级！")
		return nil
	}

	cmds := []string{update}

	for _, cmd := range cmds {
		fmt.Println(cmd)
		cmd := exec.Command("bash", "-c", cmd)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
