package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"strings"
)

//go:embed xerror.go.tpl
var tpl string

type errorInfo struct {
	Name string // Greeter

	Comment string
	Number  int32
	Msg     string
}

type errorInfos struct {
	Errors []*errorInfo
}

func (s *errorInfos) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("xerror").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}

	return buf.String()
}
