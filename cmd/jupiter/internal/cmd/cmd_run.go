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

package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cosmtrek/air/runner"
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
