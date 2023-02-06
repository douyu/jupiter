package main

import (
	"fmt"

	"github.com/douyu/jupiter/pkg"
	errorv1 "github.com/douyu/jupiter/proto/error/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

var version = pkg.JupiterVersion()

const (
	xerrorPkg          = protogen.GoImportPath("github.com/douyu/jupiter/pkg/util/xerror")
	deprecationComment = "// Deprecated: Do not use."
)

// generateFile generates a _error.pb.go file.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_xerror.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by github.com/douyu/jupiter/cmd/protoc-gen-go-xerror. DO NOT EDIT.")
	g.P("// versions:")
	g.P(fmt.Sprintf("// - protoc-gen-go-xerror %s", version))
	g.P("// - protoc             ", protocVersion(gen))
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the github.com/douyu/jupiter/cmd/protoc-gen-go-xerror package it is being compiled against.")
	g.P("var _ = new(", xerrorPkg.Ident("Err"), ")")
	g.P()

	genErrors(gen, file, g)

	return g
}

func genErrors(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	errorInfos := new(errorInfos)
	for _, enum := range file.Enums {
		for _, enumValue := range enum.Values {
			msg := proto.GetExtension(enumValue.Desc.Options(), errorv1.E_Msg)
			if msg != nil && msg.(string) != "" {
				errorInfos.Errors = append(errorInfos.Errors, &errorInfo{
					Name:    string(enumValue.Desc.Name()),
					Number:  int32(enumValue.Desc.Number()),
					Msg:     msg.(string),
					Comment: enumValue.Comments.Leading.String(),
				})
			}
		}
	}

	g.P(errorInfos.execute())
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}

	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}

	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}
