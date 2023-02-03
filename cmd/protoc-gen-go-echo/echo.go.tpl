type {{ $.InterfaceName }} interface {
{{range .MethodSet}}
	{{.Comment}}{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}
func Register{{ $.InterfaceName }}(r *v4.Echo, srv {{ $.InterfaceName }}) {
	s := _echo_{{$.Name}}{
		server: srv,
		router:     r,
	}
	s.registerService()
}

type _echo_{{$.Name}} struct{
	server {{ $.InterfaceName }}
	router *v4.Echo
}

{{range .Methods}}
{{.Comment}}func (s *_echo_{{$.Name}}) _handler_{{ .HandlerName }} (ctx v4.Context) error {
	var in {{.Request}}
	if err := ctx.Bind(&in); err != nil {
		ctx.Error(err)
		return nil
	}
	md := metadata.New(nil)
	for k, v := range ctx.Request().Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx.Request().Context(), md)
	out, err := s.server.({{ $.InterfaceName }}).{{.Name}}(newCtx, &in)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, out)
}
{{end}}

func (s *_echo_{{$.Name}}) registerService() {
{{range .Methods}}
	{{.Comment}}s.router.Add("{{.Method}}", "{{.Path}}", s._handler_{{ .HandlerName }})
{{end}}
}
