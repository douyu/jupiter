type {{ $.InterfaceName }} interface {
{{range .MethodSet}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}
func Register{{ $.InterfaceName }}(r *echo.Echo, srv {{ $.InterfaceName }}) {
	s := {{.Name}}{
		server: srv,
		router:     r,
	}
	s.RegisterService()
}

type {{$.Name}} struct{
	server {{ $.InterfaceName }}
	router *echo.Echo
}

{{range .Methods}}
func (s *{{$.Name}}) {{ .HandlerName }} (ctx echo.Context) error {
	var in {{.Request}}
	if err := xhttp.DefaultProtoBinder.Bind(&in, ctx); err != nil {
	    return xhttp.ProtoJSON(ctx, http.StatusOK, err)
	}
	md := metadata.New(nil)
	for k, v := range ctx.Request().Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx.Request().Context(), md)
	out, err := s.server.({{ $.InterfaceName }}).{{.Name}}(newCtx, &in)
	if err != nil {
		return xhttp.ProtoJSON(ctx, http.StatusOK, err)
	}

	return xhttp.ProtoJSON(ctx, http.StatusOK, out)
}
{{end}}

func (s *{{$.Name}}) RegisterService() {
{{range .Methods}}
		s.router.Add("{{.Method}}", "{{.Path}}", s.{{ .HandlerName }})
{{end}}
}
