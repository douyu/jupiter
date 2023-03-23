package main

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed fieldmask.go.tpl
var tpl string

type service struct {
	Name     string
	FullName string
	FilePath string
	Message  []*message
}

type message struct {
	In                  bool
	Out                 bool
	RequestName         string
	ResponseName        string
	UpdateInFields      []*field
	UpdateOutFields     []*field
	IdentifyFieldGoName string
	IdentifyField       string
}

type field struct {
	UnderLineName string `json:"UnderLineName"`
	DotName       string `json:"DotName"`
}

type MsgType uint32

const (
	MsgTypeIn  MsgType = 1
	MsgTypeOut MsgType = 2
)

func (s *service) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("fieldmask").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return buf.String()
}

func (msg *message) RemovePrefix() {
	for _, f := range msg.UpdateInFields {
		f.DotName = removePrefix(f.DotName)
		f.UnderLineName = toCamelCase(removePrefix(f.UnderLineName))
	}
	for _, f := range msg.UpdateOutFields {
		f.DotName = removePrefix(f.DotName)
		f.UnderLineName = toCamelCase(removePrefix(f.UnderLineName))
	}
}
func removePrefix(name string) string {
	if name == "" {
		return name
	}
	if name[0] == '.' || name[0] == '_' {
		return name[1:]
	}
	return name
}

func toCamelCase(name string) string {
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '.' || r == '_'
	})
	for i := range words {
		words[i] = cases.Title(language.English).String(words[i])
	}
	return strings.Join(words, "")
}
