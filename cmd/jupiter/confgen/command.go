// Copyright 2021 rex lv
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

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./testdata/config.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	pvv := &configStructVisitor{}
	ast.Walk(pvv, f)
}

type configStructVisitor struct {
}

func (v *configStructVisitor) Visit(n ast.Node) (w ast.Visitor) {
	tspec, ok := n.(*ast.TypeSpec)
	if !ok {
		return v
	}

	if tspec.Name.Name == "Config" {
		fmt.Printf("config=%+v\n", 111)
	}

	if tspec.Name.Name == "ProducerConfig" {
		fmt.Printf("222=%+v\n", 222)
	}

	return v

	// ttype, ok := tspec.Type.(*ast.StructType)
	// if !ok {
	// 	return v
	// }

	// for _, field := range ttype.Fields.List {
	// 	fmt.Printf("field=%+v\n", field.Tag)

	// 	stype, ok := field.Type.(*ast.SelectorExpr)
	// 	if !ok {
	// 		return v
	// 	}
	// }
	// return v
}

type configBuilderVisitor struct {
	pkgPath string
}

func (v *configBuilderVisitor) Visit(n ast.Node) (w ast.Visitor) {
	return v
}
