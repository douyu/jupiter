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
	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
)

func init() {

	tpl := `                                     
	(_)_   _ _ __ (_) |_ ___ _ __
	| | | | | '_ \| | __/ _ \ '__|
	| | |_| | |_) | | ||  __/ |
   _/ |\__,_| .__/|_|\__\___|_|
  |__/      |_|	  
							   
GoVersion: {{ .GoVersion }}
GOOS: {{ .GOOS }}
GOARCH: {{ .GOARCH }}
NumCPU: {{ .NumCPU }}
GOPATH: {{ .GOPATH }}
GOROOT: {{ .GOROOT }}
Compiler: {{ .Compiler }}
ENV: {{ .Env "GOPATH" }}
Now: {{ .Now "Monday, 2 Jan 2006" }}

`
	banner.InitString(colorable.NewColorableStdout(), true, true, tpl)
}
