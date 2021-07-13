package template

var GRPCServerTemplate = `
package {{.Package.Name}}

import (
	"context"
    {{if not .Prefix}}
    "pb/{{.Package.Name}}"
    {{else}}
    "{{.Prefix}}/pb/{{.Package.Name}}"
    {{end}}
)

type {{.Service.Name}}Server struct{}
{{range .RPC}}
func (server *{{$.Service.Name}}Server) {{.Name}}(context context.Context, request *{{$.Package.Name}}.{{.RequestType}}) (response *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	panic("implement me")
    return 
}
{{end}}
`
