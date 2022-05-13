// Copyright 2022 Douyu
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

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/urfave/cli"
	"github.com/xlab/treeprint"
)

const (
	// 存放于git上的模板module
	gitOriginModulePath = "github.com/douyu/jupiter-layout"
)

// New 生成项目
func New(c *cli.Context) error {
	return generate(c, c.String("type"))
}

// generate 生成项目
func generate(c *cli.Context, cmd string) error {
	if len(c.Args().First()) == 0 {
		return errors.New("no project name like test-go found")
	}

	// fmt.Println("cmd name:", cmd)

	dir := c.Args().First()

	goDir := filepath.Join(path.Clean(dir))

	files := make([]file, 0)

	gitFileInfos := getFileInfosByGit(cmd)
	for _, f := range gitFileInfos {
		files = append(files, *f)
	}

	cfg := config{
		Name:  generateName(dir),
		Type:  cmd,
		Dir:   dir,
		GoDir: goDir,
		Files: files,
		Comments: []string{
			// "run to compile......",
			fmt.Sprintf("Generate %s project success", cmd),
			"\ncd  " + filepath.Join(goDir),
			"\njupiter run -c cmd/exampleserver/.jupiter.toml",
			"\nEnjoy coding~~",
		},
	}

	if err := create(cfg); err != nil {
		return err
	}

	return nil
}

type file struct {
	Path string
	Tmpl []byte
}

// test-go => test
func generateName(pname string) string {
	ss := strings.Split(pname, "/")
	ss = strings.Split(ss[len(ss)-1], "-")
	return ss[0]
}

type config struct {
	// foo
	Alias string
	// api, srv, web, task, admin
	Type string
	// test-go
	Dir string
	// $GOPATH/src/test-go
	GoDir string
	// $GOPATH
	GoPath string
	// test
	Name string
	// TemplatePath
	TemplatePath string
	// Files
	Files []file

	IsFront bool
	// Comments
	Comments []string
}

// create 根据配置来生成模板对应的新项目
func create(c config) error {
	// check if dir exists
	// if _, err := os.Stat(c.GoDir); !os.IsNotExist(err) {
	// 	return fmt.Errorf("%s already exists", c.GoDir)
	// }

	fmt.Printf("Creating service with an app [exampleserver] in %s\n\n", c.GoDir)

	t := treeprint.New()

	nodes := map[string]treeprint.Tree{}
	nodes[c.GoDir] = t

	// 按ASCII字符排序，忽略大小写
	sort.Slice(c.Files, func(i, j int) bool {
		return strings.Compare(strings.ToLower(c.Files[i].Path), strings.ToLower(c.Files[j].Path)) < 0
	})

	// write the files
	for _, file := range c.Files {
		f := filepath.Join(c.GoDir, file.Path)
		dir := filepath.Dir(f)

		b, ok := nodes[dir]
		if !ok {
			d, _ := filepath.Rel(c.GoDir, dir)
			if !strings.Contains(f, "vendor") {
				// 不打印vendor下目录
				b = t.AddBranch(d)
				nodes[dir] = b
			}
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		p := filepath.Base(f)
		if !strings.Contains(f, "vendor") {
			// vendor依赖目录不打印文件
			b.AddNode(p)
		}

		// 计算相对于$GOPATH/src的路径之后
		// 替换为真实的c.Dir
		tpl := strings.ReplaceAll(
			string(file.Tmpl),
			gitOriginModulePath, //+"/"+c.Type,
			c.Dir)

		if err := write(c, f, tpl); err != nil {
			return err
		}
	}

	// print tree
	fmt.Println(t.String())

	for _, comment := range c.Comments {
		fmt.Println(comment)
	}

	return nil
}

func write(c config, file, tmpl string) error {

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(tmpl))
	return err
}

// getFileInfosByGit 从git拉取最新的模板代码 并抽象成map[相对路径]文件流
// name: 生成的项目类型
func getFileInfosByGit(name string) (
	fileInfos map[string]*file) {
	fileInfos = make(map[string]*file)
	// 存放于git的模板地址
	gitPath := "https://" + gitOriginModulePath + ".git"

	// clone最新仓库的master分支
	// 设置clone别名 避免冲突
	finalTemplateDir := "local_temp_jupiter_layout"
	fmt.Println("git", "clone", gitPath)

	cmd := exec.Command("git", "clone", gitPath, finalTemplateDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("[jupiter] Clone template git codebase failed:template_name=[%v] gitPath=[%v] err=[%v]",
			name, gitPath, err)
		return fileInfos
	}

	// 获取模板的文件流
	// io/fs为1.16新增标准库 低版本不支持
	// os.FileInfo实现了和io/fs.FileInfo相同的接口 确保go低版本可以成功编译通过
	err := filepath.Walk("./"+finalTemplateDir, func(path string, info os.FileInfo, err error) error {

		// 过滤git目录中文件
		if !info.IsDir() && !strings.Contains(path, ".git/") {

			fullPath := strings.ReplaceAll(path, "\\", "/")
			bs, err := ioutil.ReadFile(fullPath)
			if err != nil {
				fmt.Printf("[jupiter] Read file failed: fullPath=[%v] err=[%v]",
					fullPath, err)

			}

			fullPath = strings.ReplaceAll(fullPath, finalTemplateDir, "")
			fileInfos[fullPath] = &file{fullPath, bs}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// 删除clone的仓库
	err = os.RemoveAll("./" + finalTemplateDir)
	if err != nil {
		panic(err)
	}

	return fileInfos
}
