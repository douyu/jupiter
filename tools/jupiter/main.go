// Copyright 2020 Douyu
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

package main

import (
	"github.com/douyu/jupiter/tools/jupiter/new"
	"github.com/douyu/jupiter/tools/jupiter/protoc"
	"log"
	"os"

	"github.com/urfave/cli"
)

const Version = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "jupiter"
	app.Usage = "jupiter tools"
	app.Version = Version
	app.Commands = []cli.Command{
		new.Cmd,
		protoc.Cmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
