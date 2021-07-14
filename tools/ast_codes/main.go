// Copyright 2020 Douyu
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
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/douyu/jupiter/pkg/xlog"

	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/fatih/structtag"
	"github.com/flosch/pongo2"

	// "github.com/pelletier/go-toml"
	"github.com/davecgh/go-spew/spew"
	"github.com/douyu/pkg/util/xcast"
)

/*
usage: go run main.go --dir=../../
*/
var packages = []string{
	"pkg/client/etcdv3",
	"pkg/client/grpc",
	"pkg/client/redis",
	"pkg/client/rocketmq",
	"pkg/server/xecho",
	"pkg/server/xgrpc",
	"pkg/store/gorm",
	"./",
}

var configMap = map[string]string{
	"client/etcdv3": "jupiter.etcdv3",
	//"client/grpc":   "jupiter.grpc",
}

var fieldMap = map[string]string{
	"FieldAID":         "aid",
	"FieldName":        "name",
	"FieldAddr":        "addr",
	"FieldEvent":       "event",
	"FieldType":        "tp",
	"FieldModule":      "mod",
	"FieldCounter":     "tid",
	"FieldBegTime":     "beg",
	"FieldEndTime":     "end",
	"FieldCost":        "cost",
	"FieldError":       "err",
	"FieldStringError": "err",
	"FieldHost":        "host",
	"FieldCode":        "code",
	"FieldReqAID":      "reqAid",
	"FieldStdMethod":   "stdMthd",
	"FieldTID":         "tid",
	"FieldIP":          "ip",
	"FieldReqHost":     "reqHost",
	"FieldExtMessage":  "ext",
	"FieldRestrictExt": "ext",
	"FieldStack":       "stack",
	"FieldMethod":      "mthd",
}

func init() {
	flag.Register(
		&flag.StringFlag{
			Name:    "dir",
			Usage:   "--dir ",
			Default: ".",
		},
	)
}

/*
解析各个包，并生成文档
*/

func main() {
	_ = flag.Parse()
	spew.Dump("")

	xlog.Info("ast codes start")
	for _, pkg := range packages {
		fset := token.NewFileSet()
		fpkgs, err := parser.ParseDir(fset, "../../"+pkg, nil, parser.ParseComments)
		if err != nil {
			xlog.Panic("pkg parse panic", xlog.FieldErr(err))
		}

		for _, f := range fpkgs {
			pvv := &MsgInfoVistor{
				data:  map[string]string{},
				nodes: make([]*Node, 0),
				name:  pkg,
				fset:  fset,
			}
			xlog.Info("pkg wark", xlog.FieldMod(pkg))
			ast.Walk(pvv, f)
		}
	}

	var data = make(map[string]map[string][]map[string]interface{})
	for mod, ns := range nodes {
		if data[mod] == nil {
			data[mod] = make(map[string][]map[string]interface{})
		}

		xlog.Info("nodes info", xlog.FieldMod(mod), xlog.FieldValueAny(ns))

		for _, n := range ns {
			if data[mod][n.Level] == nil {
				data[mod][n.Level] = make([]map[string]interface{}, 0)
			}

			var fields = []string{"mod"}
			for _, field := range n.Fields {
				if fm, ok := fieldMap[field]; ok {
					fields = append(fields, fm)
				} else {
					fields = append(fields, field)
				}
			}

			filepath := strings.TrimLeft(n.Position.Filename, flag.String("dir"))
			data[mod][n.Level] = append(data[mod][n.Level], map[string]interface{}{
				"message":  n.Message,
				"fields":   strings.Join(fields, ","),
				"position": "https://github.com/douyu/jupiter/blob/master/" + filepath + "#L" + strconv.Itoa(n.Position.Line),
			})
		}
	}

	/*
		## cache/bigmap
			- panic: reenetry (addr, mod, err)
			- panic: hello (addr, mod, err)
			- warn: world (addr, mod, err)
	*/

	var codeMDs = []map[string]interface{}{}
	for mod, lvs := range data {
		var mdData = []map[string]interface{}{}
		for lv, ele := range lvs {
			for _, e := range ele {
				mdData = append(mdData, map[string]interface{}{
					"level":    lv,
					"message":  e["message"],
					"fields":   e["fields"],
					"position": e["position"],
				})
			}
		}
		codeMDs = append(codeMDs, map[string]interface{}{
			"name":   mod,
			"fields": mdData,
		})
	}
	//
	// xdebug.PrettyJsonPrint("configs", configDescs)
	// xdebug.PrettyJsonPrint("data", data)

	// tree, err := toml.LoadFile(flag.String("dir") + "Document/规范/配置示例.toml")
	// if err != nil {
	// 	panic(err)
	// }

	var configMDs = []map[string]interface{}{}
	for _, desc := range configDescs {
		var mdData = map[string]interface{}{}
		if _, ok := configMap[desc.Package]; ok {
			var mdfields = make([]MDConfig, 0)
			for _, field := range desc.Fields {
				var mdfield = MDConfig{
					Name: field.Tag,
					Doc:  field.Doc,
				}
				switch field.Type {
				case "bool":
					// tree.SetWithComment(prefix+".default."+field.Tag, field.Doc, false, xcast.ToBool(field.Default))
					mdfield.Default = xcast.ToBool(field.Default)
				case "int", "int32", "int64", "int8", "uint", "uint32", "uint64":
					// tree.SetWithComment(prefix+".default."+field.Tag, field.Doc, false, xcast.ToInt64(field.Default))
					mdfield.Default = xcast.ToInt64(field.Default)
				case "time.Duration":
					// tree.SetWithComment(prefix+".default."+field.Tag, field.Doc, false, xcast.ToString(field.Default))
					mdfield.Default = xcast.ToString(field.Default)
				}
				mdfields = append(mdfields, mdfield)
			}

			filepath := strings.TrimLeft(desc.Position.Filename, flag.String("dir"))
			if desc.Name != "" {
				// name := strings.Replace(desc.Package, "/", ".", -1)
				mdData["pkg"] = desc.Package
				mdData["name"] = desc.Name
				mdData["fields"] = mdfields
				mdData["position"] = "https://github.com/douyu/jupiter/blob/master/" + filepath + "#L" + strconv.Itoa(desc.Position.Line)
			}

			fmt.Printf("mdData = %+v\n", mdData["position"])
		}
		configMDs = append(configMDs, mdData)
	}

	// xdebug.PrettyJsonPrint("md", configMDs)

	// file, err := os.OpenFile(flag.String("dir")+"Document/规范/配置示例.toml", os.O_RDWR, 0777)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// defer file.Close()
	// _, _ = tree.WriteTo(file)

	// _ = ioutil.WriteFile(flag.String("dir")+"Document/Codes.json", xstring.PrettyJSONBytes(data), 0777)
	// _ = ioutil.WriteFile(flag.String("dir")+"Document/config.json", xstring.PrettyJSONBytes(configDescs), 0777)

	{
		f, err := os.Create("./配置说明.md")
		if err != nil {
			panic(err)
		}

		tpl, err := pongo2.FromString(configMDTemplate)
		if err != nil {
			panic(err)
		}

		if err := tpl.ExecuteWriter(pongo2.Context(map[string]interface{}{"configs": configMDs}), f); err != nil {
			panic(err)
		}
	}

	{
		f, err := os.Create("./错误日志埋点.md")
		if err != nil {
			panic(err)
		}

		tpl, err := pongo2.FromString(codesStatusTableTempldate)
		if err != nil {
			panic(err)
		}

		xdebug.PrintObject("codes", codeMDs)
		if err := tpl.ExecuteWriter(pongo2.Context(map[string]interface{}{"packages": codeMDs}), f); err != nil {
			panic(err)
		}
	}
}

// MDConfig ...
type MDConfig struct {
	Name    string      `json:"name" toml:"name"`
	Doc     string      `json:"doc" toml:"doc"`
	Default interface{} `json:"default" toml:"default"`
}

var configMDTemplate = `
# mienrva 配置结构列表

> 本文档是由github.com/douyu/jupiter/script/ast_codes自动生成   
> 配置格式有变更时，请执行: cd script/ast_codes; go run main.go --dir=../../; cd -  

{% for config in configs %}
## {{config.pkg}}
### {{config.name}}  
> 点我: {{config.position}}
{% for field in config.fields %}
- {{field.Name}}: {{field.Doc}}  
	(default: {{field.Default}})
{% endfor %}
{% endfor %}
`

var codesStatusTableTempldate = `
{% for pkg in packages %}
## {{pkg.name}} 
|  错误 | 级别 | 行数 |
|:--------------|:-----|:-------------------|
{% for field in pkg.fields %}| {{field.message}} | {{field.level}}|[代码地址]({{field.position}})|
{% endfor %}
{% endfor %}

`

// MsgInfoVistor ...
type MsgInfoVistor struct {
	name  string
	data  map[string]string
	pkg   string
	nodes []*Node
	fset  *token.FileSet
}

// Node ...
type Node struct {
	Module   string         `json:"-" toml:"-"`
	Level    string         `json:"level" toml:"level"`
	Message  string         `json:"message" toml:"message"`
	Fields   []string       `json:"fields" toml:"fields"`
	Position token.Position `json:"position" toml:"position"` // 在代码中的位置
}

// String ...
func (node *Node) String() string {
	var buff bytes.Buffer
	buff.WriteString(node.Module + ":")
	buff.WriteString("\t" + node.Level)
	buff.WriteString("\t" + node.Message)
	buff.WriteString("\t" + strings.Join(node.Fields, ","))
	return buff.String()
}

var nodes = make(map[string][]*Node)

// ConfigDesc ...
type ConfigDesc struct {
	Package  string
	Name     string
	Fields   map[string]*FieldDesc
	Position token.Position
}

// FieldDesc ...
type FieldDesc struct {
	Name    string
	Tag     string
	Doc     string
	Type    string
	Default interface{}
}

var configDescs = make([]ConfigDesc, 0)

// 将节点数据解析到全局变量里
func (pvv *MsgInfoVistor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	// 解析调用信息
	/**
	(*ast.CallExpr)(0xc000033340)({
	 Fun: (*ast.SelectorExpr)(0xc00000fb80)({
	  X: (*ast.Ident)(0xc00000fb40)(context),
	  Sel: (*ast.Ident)(0xc00000fb60)(WithTimeout)
	 }),
	 Lparen: (token.Pos) 553,
	 Args: ([]ast.Expr) (len=2 cap=2) {
	  (*ast.CallExpr)(0xc000033300)({
	   Fun: (*ast.SelectorExpr)(0xc00000fbe0)({
		X: (*ast.Ident)(0xc00000fba0)(context),
		Sel: (*ast.Ident)(0xc00000fbc0)(Background)
	   }),
	   Lparen: (token.Pos) 572,
	   Args: ([]ast.Expr) <nil>,
	   Ellipsis: (token.Pos) 0,
	   Rparen: (token.Pos) 573
	  }),
	  (*ast.Ident)(0xc00000fc00)(timeout)
	 },
	 Ellipsis: (token.Pos) 0,
	 Rparen: (token.Pos) 583
	})
	*/
	case *ast.CallExpr:
		switch fun := t.Fun.(type) {
		case *ast.SelectorExpr:
			switch xType := fun.X.(type) {
			// 日志都是使用结构体的.logger进行使用
			case *ast.SelectorExpr:
				recv := xType.Sel
				// 如果为日志，记录到node里
				if recv.Name == "logger" {
					// 去除With方法
					if fun.Sel.Name == "With" {
						break
					}
					var node = &Node{
						Module:   pvv.name,
						Level:    fun.Sel.Name,
						Fields:   make([]string, 0),
						Position: pvv.fset.Position(fun.Pos()),
					}
					for _, arg := range t.Args {
						if lit, ok := arg.(*ast.BasicLit); ok {
							node.Message = strings.Trim(lit.Value, "\"")
						}

						if call, ok := arg.(*ast.CallExpr); ok {
							if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
								node.Fields = append(node.Fields, sel.Sel.Name)
							}
						}
					}

					if nodes[pvv.name] == nil {
						nodes[pvv.name] = make([]*Node, 0)
					}
					nodes[pvv.name] = append(nodes[pvv.name], node)
					xlog.Info("logger info", xlog.FieldMod(pvv.name), xlog.String("level", fun.Sel.Name), xlog.String("msg", node.Message), xlog.Any("pos", pvv.fset.Position(fun.Pos())))
				}

			}
		}
	case *ast.GenDecl:
		for _, spec := range t.Specs {
			// 解析配置结构
			ts, ok := spec.(*ast.TypeSpec)
			if ok {
				if strings.HasSuffix(ts.Name.Name, "config") {
					configDesc := ConfigDesc{
						Package:  pvv.name,
						Name:     ts.Name.String(),
						Fields:   make(map[string]*FieldDesc),
						Position: pvv.fset.Position(ts.Pos()),
					}
					if st, ok := ts.Type.(*ast.StructType); ok {
						for _, field := range st.Fields.List {
							if field.Tag == nil || field.Tag.Value == "" {
								continue
							}
							var vv string
							switch value := field.Type.(type) {
							case *ast.Ident:
								vv = value.Name
							case *ast.SelectorExpr:
								if fmt.Sprint(value.X) == "time" {
									vv = fmt.Sprintf("%v.%v", value.X, value.Sel.Name)
								}
							}
							value := strings.Trim(field.Tag.Value, "`")
							tags, err := structtag.Parse(string(value))
							if err != nil {
								panic(err)
							}
							if tags == nil {
								continue
							}

							tag, err := tags.Get("toml")
							if err != nil {
								// todo
								//	panic(err)
							} else {
								configDesc.Fields[field.Names[0].Name] = &FieldDesc{
									Name: field.Names[0].Name,
									Tag:  tag.Value(),
									Doc:  field.Doc.Text(),
									Type: vv,
								}
							}

						}
					}

					configDescs = append(configDescs, configDesc)
				}
			}
		}
	case *ast.FuncDecl:
		// 解析默认配置
		if strings.HasSuffix(t.Name.Name, "config") && strings.HasPrefix(t.Name.Name, "Default") {
			for _, stmt := range t.Body.List {
				if ret, ok := stmt.(*ast.ReturnStmt); ok {
					for _, expr := range ret.Results {
						if lit, ok := expr.(*ast.CompositeLit); ok {
							for _, elt := range lit.Elts {
								if kv, ok := elt.(*ast.KeyValueExpr); ok {
									switch value := kv.Value.(type) {
									case *ast.CallExpr:
										if ss, ok := value.Fun.(*ast.SelectorExpr); ok {
											recv, ok := ss.X.(*ast.Ident)
											if ok {
												if recv.Name == "xtime" {
													if len(value.Args) > 0 {
														// fmt.Printf("value.Args = %+v\n", value.Args[0])
														vv := value.Args[0].(*ast.BasicLit).Value
														for _, cfg := range configDescs {
															if cfg.Package == pvv.name {
																for _, field := range cfg.Fields {
																	if key, ok := kv.Key.(*ast.Ident); ok {
																		if field.Name == key.Name {
																			field.Default = strings.Trim(vv, "\"")
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									case *ast.BinaryExpr:
										var vv string
										if ss, ok := value.X.(*ast.SelectorExpr); ok {
											if fmt.Sprint(ss.X) == "time" {
												vv = fmt.Sprintf("%v.%v * %v", ss.X, ss.Sel.Name, value.Y.(*ast.BasicLit).Value)
											}
										}

										if ss, ok := value.Y.(*ast.SelectorExpr); ok {
											if fmt.Sprint(ss.X) == "time" {
												vv = fmt.Sprintf("%v.%v * %v", ss.X, ss.Sel.Name, value.X.(*ast.BasicLit).Value)
											}
										}

										for _, cfg := range configDescs {
											if cfg.Package == pvv.name {
												for _, field := range cfg.Fields {
													if key, ok := kv.Key.(*ast.Ident); ok {
														if field.Name == key.Name {
															field.Default = vv
														}
													}
												}
											}
										}
									case *ast.Ident:
										for _, cfg := range configDescs {
											if cfg.Package == pvv.name {
												for _, field := range cfg.Fields {
													if key, ok := kv.Key.(*ast.Ident); ok {
														if field.Name == key.Name {
															field.Default = value.Name
														}
													}
												}
											}
										}
									case *ast.BasicLit:
										for _, cfg := range configDescs {
											if cfg.Package == pvv.name {
												for _, field := range cfg.Fields {
													if key, ok := kv.Key.(*ast.Ident); ok {
														if field.Name == key.Name {
															field.Default = value.Value
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	default:
	}
	return pvv
}
