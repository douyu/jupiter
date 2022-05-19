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
	return generate(c, c.String("type"), c.Bool("clean"))
}

// generate 生成项目
func generate(c *cli.Context, cmd string, clean bool) error {
	if len(c.Args().First()) == 0 {
		return errors.New("no project name like test-go found")
	}

	// fmt.Println("cmd name:", cmd)

	dir := c.Args().First()

	goDir := filepath.Join(path.Clean(dir))

	files := make([]file, 0)

	gitFileInfos := getFileInfosByGit(cmd, clean)
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
			"\ncd " + goDir,
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
			gitOriginModulePath, // +"/"+c.Type,
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
func getFileInfosByGit(name string, clean bool) (fileInfos map[string]*file) {
	if clean {
		if err := cleanTempLayout(); err != nil {
			panic(err)
		}
	}
	// 设置clone别名 避免冲突
	// 查看临时文件之中是否已经存在该文件夹
	tempPath := filepath.Join(os.TempDir(), "local_temp_jupiter_layout")
	// os.Stat 获取文件信息
	_, err := os.Stat(tempPath)
	if os.IsNotExist(err) {
		// 存放于git的模板地址
		gitPath := "https://" + gitOriginModulePath + ".git"

		fmt.Println("git", "clone", gitPath)

		// clone最新仓库的master分支
		// 不存在则拉取模板
		cmd := exec.Command("git", "clone", gitPath, tempPath, "-b", "main", "--depth=1")
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	} else if os.IsExist(err) || err == nil {
		// 	判断是否需要刷新模板信息
		// todo ... 后面有时间再加上
	} else {
		// 这里的错误，是说明出现了未知的错误，应该抛出
		panic(err)
	}

	fileInfos = make(map[string]*file)
	// 获取模板的文件流
	// io/fs为1.16新增标准库 低版本不支持
	// os.FileInfo实现了和io/fs.FileInfo相同的接口 确保go低版本可以成功编译通过
	err = filepath.Walk(tempPath, func(path string, info os.FileInfo, err error) error {
		// 过滤git目录中文件
		if !info.IsDir() && !strings.Contains(strings.ReplaceAll(path, "\\", "/"), ".git/") {
			bs, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("[jupiter] Read file failed: fullPath=[%v] err=[%v]", path, err)
			}

			fullPath := strings.ReplaceAll(path, tempPath, "")
			fileInfos[fullPath] = &file{fullPath, bs}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return fileInfos
}
