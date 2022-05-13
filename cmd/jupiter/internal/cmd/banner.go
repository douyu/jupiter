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
