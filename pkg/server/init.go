// Copyright 2021 rex lv
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

package server

import (
	"log"

	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/util/xnet"
)

func init() {
	defaultHost, _, err := xnet.GetLocalMainIP()
	if err != nil {
		defaultHost = "127.0.0.1"
	}

	flag.Register(&flag.StringFlag{
		Name:    "host",
		Usage:   "--host, print host",
		Default: defaultHost,
		Action: func(name string, fs *flag.FlagSet) {
			log.Printf("host flag: %v", fs.String(name))
		},
	})
}
