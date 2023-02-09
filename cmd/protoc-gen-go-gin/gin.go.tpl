type {{ $.InterfaceName }} interface {
{{range .MethodSet}}
	{{.Comment}}{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}
func Register{{ $.InterfaceName }}(r *gin.Engine, srv {{ $.InterfaceName }}) {
	s := _gin_{{$.Name}}{
		server: srv,
		router:     r,
	}
	s.registerService()
}

type _gin_{{$.Name}} struct{
	server {{ $.InterfaceName }}
	router *gin.Engine
}

{{range .Methods}}
{{.Comment}}func (s *_gin_{{$.Name}}) _gin_handler_{{ .HandlerName }} (ctx *gin.Context) {
	var in {{.Request}}
	if err := ctx.ShouldBindUri(&in); err != nil {
		ctx.Error(err)
		return
	}
	if err := ctx.ShouldBind(&in); err != nil {
		ctx.Error(err)
		return
	}
	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx.Request.Context(), md)
	out, err := s.server.({{ $.InterfaceName }}).{{.Name}}(newCtx, &in)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, out)
}
{{end}}

func (s *_gin_{{$.Name}}) registerService() {
{{range .Methods}}
	{{.Comment}}s.router.Handle("{{.Method}}", "{{.Path}}", s._gin_handler_{{ .HandlerName }})
{{end}}
}
