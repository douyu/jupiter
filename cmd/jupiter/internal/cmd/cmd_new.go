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
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/xlab/treeprint"
)

const (
	// 存放于git上的模板module
	gitOriginModulePath = "github.com/douyu/jupiter-layout"
	gitOriginModuleName = "douyu/jupiter-layout"
	// 设置clone别名 避免冲突
	localLayoutName = "local_temp_jupiter_layout"
	// 本地文件锁，修改时间用于检查更新
	localLayoutNameLock = "local_temp_jupiter_layout.lock"
	oneDayUnix          = 24 * int64(time.Hour)
)

var (
	// 项目模板存在的位置
	globalLayoutPath = filepath.Join(os.TempDir(), localLayoutName)
	// 项目时间检查的文件所在位置
	globalLayoutLockPath = filepath.Join(os.TempDir(), localLayoutNameLock)
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
	// 判断是否需要拉取模板
	var needCloneFile bool

	// 查看临时文件之中是否已经存在该文件夹
	// os.Stat 获取文件信息
	_, err := os.Stat(globalLayoutPath)
	if os.IsNotExist(err) {
		needCloneFile = true
	} else if err != nil {
		// 这里的错误，是说明出现了未知的错误，应该抛出
		panic(err)
	}

	// 判断是否需要刷新模板信息
	// 存在文件才检查更新
	if (err == nil && checkUpgrade()) || clean {
		if err := cleanTempLayout(); err != nil {
			panic(err)
		}
		needCloneFile = true
	}

	if needCloneFile {
		// 存放于git的模板地址
		gitPath := "https://" + gitOriginModulePath + ".git"

		fmt.Println("git", "clone", gitPath)

		// clone最新仓库的master分支
		// 不存在则拉取模板
		cmd := exec.Command("git", "clone", gitPath, globalLayoutPath, "-b", "main", "--depth=1")
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	fileInfos = make(map[string]*file)
	// 获取模板的文件流
	// io/fs为1.16新增标准库 低版本不支持
	// os.FileInfo实现了和io/fs.FileInfo相同的接口 确保go低版本可以成功编译通过
	err = filepath.Walk(globalLayoutPath, func(path string, info os.FileInfo, err error) error {
		// 过滤git目录中文件
		if !info.IsDir() && !strings.Contains(strings.ReplaceAll(path, "\\", "/"), ".git/") {
			bs, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("[jupiter] Read file failed: fullPath=[%v] err=[%v]", path, err)
			}

			fullPath := strings.ReplaceAll(path, globalLayoutPath, "")
			fileInfos[fullPath] = &file{fullPath, bs}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return fileInfos
}

// checkUpgrade 通过 https://api.github.com/repos/xxx/commits/branch 获取最后一次提交sha,判断本地提交是否一致，不一致则需要更新
func checkUpgrade() bool {
	color.Green("check upgrade ....")

	// 检查今天是否已经检查过更新
	if !checkDays() {
		return false
	}

	// 获取远端最后一次提交的 SHA
	remoteLastSha, err := getRemoteLastCommitSha()
	if err != nil {
		return false
	}

	// 和本地对比，如果一致，那么不需要更新
	if getLocalLastSha() == remoteLastSha {
		return false
	}

	// 由用户选择是否更新
	return userSelectUpgrade()
}

// checkDays
// 	存在文件，首先判断时间是否需要进行更新
// 	不存在文件肯定是需要创建文件，并确定需要更新
//  存在文件但是时间不是同一天，那么也需要进行更新
//  需要更新的同时，更新文件的修改时间
func checkDays() bool {
	fileInfo, err := os.Stat(globalLayoutLockPath)
	if err == nil && fileInfo.ModTime().Unix()/oneDayUnix == time.Now().Unix()/oneDayUnix {
		return false
	} else if os.IsNotExist(err) {
	} else if err != nil {
		// 这里的错误，是说明出现了未知的错误，应该抛出
		panic(err)
	}

	f, err := os.Create(globalLayoutLockPath)
	if err != nil {
		panic(err)
	}
	_ = f.Close()

	return true
}

func getRemoteLastCommitSha() (string, error) {
	request, _ := http.NewRequest("GET", "https://api.github.com/repos/"+gitOriginModuleName+"/commits/main", nil)
	request.Header.Set("Accept", "application/vnd.github.v3.sha")

	httpCli := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := httpCli.Do(request)
	if err != nil {
		color.Red("check remote upgrade filed:%v", err)
		return "", err
	}

	defer resp.Body.Close()
	remoteLastSha, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("check remote upgrade filed:%v", err)
		return "", err
	}

	return strings.TrimSpace(string(remoteLastSha)), nil
}

// getLocalLastSha 获取本地最后一次的提交sha
func getLocalLastSha() string {
	cmd := exec.Command("git", "rev-parse", "main")
	cmd.Dir = globalLayoutPath

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	return strings.TrimSpace(out.String())
}

func userSelectUpgrade() bool {
	var upgrade string
	for {
		color.Blue(`the current template has been updated, please select whether it needs to be updated?`)
		fmt.Print("yes/no:")

		_, err := fmt.Scanln(&upgrade)
		if err != nil {
			panic(err)
		}

		switch strings.TrimSpace(upgrade) {
		case "yes", "y", "Y":
			return true
		case "no", "n", "N":
			return false
		default:
			color.Red("you chose wrong please choose again")
		}
	}
}
