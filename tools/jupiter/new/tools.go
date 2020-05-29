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


package new

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/douyu/jupiter/pkg/util/xcolor"
	"github.com/douyu/jupiter/pkg/util/xregexp"
	"github.com/gobuffalo/packr/v2"
	"github.com/urfave/cli"
)

// CreateProject create a template project for Jupiter
func CreateProject(cli *cli.Context) (err error) {
	newArgs := cli.Args()
	if len(newArgs) <= 0 {
		fmt.Println(xcolor.Red("Command line new execution error, please use jupiter new -h for details"))
		return
	}
	name := newArgs[0]
	if name == "" {
		Project.Name = DefaultProjectName
	} else {
		Project.Name = name
	}
	if Project.Path != "" {
		if Project.Path, err = filepath.Abs(Project.Path); err != nil {
			return
		}
		Project.Path = filepath.Join(Project.Path, Project.Name)
	} else {
		pwd, _ := os.Getwd()
		Project.Path = filepath.Join(pwd, Project.Name)
	}
	modPath := getModPath(Project.Path)
	Project.ModPrefix = modPath
	if err = doCreateProject(); err != nil {
		return
	}
	fmt.Println(xcolor.Greenf("Project dir:", Project.Path))
	fmt.Println(xcolor.Green("Project created successfully"))
	return
}

func getModPath(projectPath string) (modPath string) {
	dir := filepath.Dir(projectPath)
	for {
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				content, _ := ioutil.ReadFile(filepath.Join(dir, "go.mod"))
				mod := xregexp.RegexpReplace(`module\s+(?P<name>[\S]+)`, string(content), "$name")
				name := strings.TrimPrefix(filepath.Dir(projectPath), dir)
				name = strings.TrimPrefix(name, string(os.PathSeparator))
				if name == "" {
					return fmt.Sprintf("%s/", mod)
				}
				return fmt.Sprintf("%s/%s/", mod, name)
			}
			parent := filepath.Dir(dir)
			if dir == parent {
				return ""
			}
			dir = parent
		}
	}

}
func doCreateProject() (err error) {
	box := packr.New("all", "./templates")
	if err = os.MkdirAll(Project.Path, 0755); err != nil {
		return
	}
	for _, name := range box.List() {
		if Project.ModPrefix != "" && name == "go.mod.tmpl" {
			continue
		}
		tmpl, _ := box.FindString(name)
		i := strings.LastIndex(name, string(os.PathSeparator))
		if i > 0 {
			dir := name[:i]
			if err = os.MkdirAll(filepath.Join(Project.Path, dir), 0755); err != nil {
				return
			}
		}
		if strings.HasSuffix(name, ".tmpl") {
			name = strings.TrimSuffix(name, ".tmpl")
		}
		if err = doWriteFile(filepath.Join(Project.Path, name), tmpl); err != nil {
			return
		}
	}

	return
}

func doWriteFile(path, tmpl string) (err error) {
	data, err := parseTmpl(tmpl)
	if err != nil {
		return
	}
	fmt.Println(xcolor.Greenf("File generated----------------------->", path))
	return ioutil.WriteFile(path, data, 0644)
}

func parseTmpl(tmpl string) ([]byte, error) {
	tmp, err := template.New("").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = tmp.Execute(&buf, Project); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
