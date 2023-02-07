{{range .Errors}}
{{.Comment}}var XERROR_{{ .Name }} = xerror.New({{ .Number }}, "{{ .Msg }}")
{{end}}
