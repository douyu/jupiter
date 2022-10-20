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
	gitOriginModuleName = "douyu/jupiter-layout"

	oneDayUnix = 24 * 60 * 60
)

func getGlobalLayoutPath(path string) string {
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "-", "_")
	return filepath.Join(os.TempDir(), path)
}

func getGlobalLayoutLockPath(path string) string {
	return getGlobalLayoutPath(path) + ".lock"
}

// New 生成项目
func New(c *cli.Context) error {
	return generate(c, c.String("remote"))
}

// generate 生成项目
func generate(c *cli.Context, remote string) error {
	if len(c.Args().First()) == 0 {
		return errors.New("no project name like test-go found")
	}

	// fmt.Println("cmd name:", cmd)

	dir := c.Args().First()

	goDir := filepath.Join(path.Clean(dir))

	files := make([]file, 0)

	gitFileInfos := getFileInfosByGit(c, remote)
	for _, f := range gitFileInfos {
		files = append(files, *f)
	}

	cfg := config{
		Name:   generateName(dir),
		Remote: remote,
		Dir:    dir,
		GoDir:  goDir,
		Files:  files,
		Comments: []string{
			// "run to compile......",
			fmt.Sprintf("Generate %s project success", remote),
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
	// remote template
	Remote string
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
				fmt.Println(d)
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
			c.Remote, // +"/"+c.Type,
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
func getFileInfosByGit(c *cli.Context, gitPath string) (fileInfos map[string]*file) {
	// 查看临时文件之中是否已经存在该文件夹
	// os.Stat 获取文件信息
	_, err := os.Stat(getGlobalLayoutPath(gitPath))
	if os.IsNotExist(err) {
		// 不存在，拉取对应的仓库
		cloneGitRepo(gitPath)
	} else if err != nil {
		// 这里的错误，是说明出现了未知的错误，应该抛出
		panic(err)
	}

	// 判断是否需要刷新模板信息
	// 存在文件才检查更新
	if err == nil && checkUpgrade(c, gitPath) {
		pullGitRepo(gitPath)
	}

	fileInfos = make(map[string]*file)
	// 获取模板的文件流
	// io/fs为1.16新增标准库 低版本不支持
	// os.FileInfo实现了和io/fs.FileInfo相同的接口 确保go低版本可以成功编译通过
	err = filepath.Walk(getGlobalLayoutPath(gitPath), func(path string, info os.FileInfo, err error) error {
		// 过滤git目录中文件
		if !info.IsDir() && !strings.Contains(strings.ReplaceAll(path, "\\", "/"), ".git/") {
			bs, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("[jupiter] Read file failed: fullPath=[%v] err=[%v]", path, err)
			}

			fullPath := strings.ReplaceAll(path, getGlobalLayoutPath(gitPath), "")
			fileInfos[fullPath] = &file{fullPath, bs}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return fileInfos
}

// 拉取 jupiter-layout 仓库
func cloneGitRepo(path string) {
	// 存放于git的模板地址
	gitPath := "https://" + path + ".git"

	fmt.Println("git", "clone", gitPath)

	// clone最新仓库的master分支
	// 不存在则拉取模板
	var stdErr bytes.Buffer
	cmd := exec.Command("git", "clone", gitPath, getGlobalLayoutPath(path), "-b", "main", "--depth=1")
	cmd.Stderr = &stdErr
	if err := cmd.Run(); err != nil {
		panic(stdErr.String())
	}
}

// 更新 jupiter-layout 仓库
func pullGitRepo(path string) {
	fmt.Println("git", "pull", path)

	// pull 最新仓库的master分支
	// 不存在则拉取模板
	var stdErr bytes.Buffer
	cmd := exec.Command("git", "pull")
	cmd.Dir = getGlobalLayoutPath(path)
	cmd.Stderr = &stdErr
	if err := cmd.Run(); err != nil {
		panic(stdErr.String())
	}
}

// checkUpgrade 通过 https://api.github.com/repos/xxx/commits/branch 获取最后一次提交sha,判断本地提交是否一致，不一致则需要更新
func checkUpgrade(c *cli.Context, path string) bool {
	if c.Bool("upgrade") {
		return true
	}

	color.Green("check upgrade (%s) ...", path)

	checkGitCorrectness(path)

	// 检查今天是否已经检查过更新
	if !checkDays(path) {
		return false
	}

	// 获取远端最后一次提交的 SHA
	remoteLastSha, err := getRemoteLastCommitSha()
	if err != nil {
		return false
	}

	// 和本地对比，如果一致，那么不需要更新
	if getLocalLastSha(path) == remoteLastSha {
		return false
	}

	// 由用户选择是否更新
	return userSelectUpgrade()
}

// 检查模板的正确性
func checkGitCorrectness(path string) {
	cmd := exec.Command("git", "status")
	cmd.Dir = getGlobalLayoutPath(path)

	var out bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr
	err := cmd.Run()
	if err == nil {
		return
	}

	if !strings.Contains(stdErr.String(), "not a git repository") {
		panic(stdErr.String())
	}

	color.Red("jupiter-layout is broken pull it again")

	err = cleanTempLayout(path)
	if err != nil {
		panic(err)
	}

	err = cleanTempLayoutLock(path)
	if err != nil {
		panic(err)
	}

	cloneGitRepo(path)

	createTempLock(path)
}

// checkDays
// 	存在文件，首先判断时间是否需要进行更新
// 	不存在文件肯定是需要创建文件，并确定需要更新
//  存在文件但是时间不是同一天，那么也需要进行更新
//  需要更新的同时，更新文件的修改时间
func checkDays(path string) bool {
	fileInfo, err := os.Stat(getGlobalLayoutLockPath(path))
	if err == nil && fileInfo.ModTime().Unix()/oneDayUnix == time.Now().Unix()/oneDayUnix {
		return false
	} else if os.IsNotExist(err) {
	} else if err != nil {
		// 这里的错误，是说明出现了未知的错误，应该抛出
		panic(err)
	}

	createTempLock(path)

	return true
}

func createTempLock(path string) {
	f, err := os.Create(getGlobalLayoutLockPath(path))
	if err != nil {
		panic(err)
	}
	_ = f.Close()
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
func getLocalLastSha(path string) string {
	cmd := exec.Command("git", "rev-parse", "main")
	cmd.Dir = getGlobalLayoutPath(path)

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
