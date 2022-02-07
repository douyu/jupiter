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

package xfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// fileInfo describes a configuration file and is returned by fileStat.
type FileInfo struct {
	Uid  uint32
	Gid  uint32
	Mode os.FileMode
	Md5  string
}

// Exists return weather file existed
func Exists(fpath string) bool {
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return false
	}

	return true
}

// ListFiles returns all file names in `dir`
func ListFiles(dir string, ext string) []string {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	var ret []string
	for _, fp := range fs {
		if fp.IsDir() {
			continue
		}

		if ext != "" && filepath.Ext(fp.Name()) != ext {
			continue
		}

		ret = append(ret, dir+"/"+fp.Name())
	}

	return ret
}

// IsFileChanged reports whether src and dest config files are equal.
// Two config files are equal when they have the same file contents and
// Unix permissions. The owner, group, and mode must match.
// Returns false in other cases.
func IsFileChanged(src, dest string) (bool, error) {
	if !Exists(dest) {
		return true, nil
	}
	d, err := FileStat(dest)
	if err != nil {
		return true, err
	}
	s, err := FileStat(src)
	if err != nil {
		return true, err
	}

	if d.Uid != s.Uid || d.Gid != s.Gid || d.Mode != s.Mode || d.Md5 != s.Md5 {
		return true, nil
	}
	return false, nil
}

// IsDirectory ...
func IsDirectory(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}

// RecursiveFilesLookup ...
func RecursiveFilesLookup(root string, pattern string) ([]string, error) {
	return recursiveLookup(root, pattern, false)
}

// RecursiveDirsLookup ...
func RecursiveDirsLookup(root string, pattern string) ([]string, error) {
	return recursiveLookup(root, pattern, true)
}

func recursiveLookup(root string, pattern string, dirsLookup bool) ([]string, error) {
	var result []string
	isDir, err := IsDirectory(root)
	if err != nil {
		return nil, err
	}
	if isDir {
		err := filepath.Walk(root, func(root string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			match, err := filepath.Match(pattern, f.Name())
			if err != nil {
				return err
			}
			if match {
				isDir, err := IsDirectory(root)
				if err != nil {
					return err
				}
				if isDir && dirsLookup {
					result = append(result, root)
				} else if !isDir && !dirsLookup {
					result = append(result, root)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if !dirsLookup {
			result = append(result, root)
		}
	}
	return result, nil
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	if runtime.GOOS == "windows" {
		dirctory = strings.Replace(dirctory, "\\", "/", -1)
	}
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

// CheckAndGetParentDir ...
func CheckAndGetParentDir(path string) string {
	// check path is the directory
	isDir, err := IsDirectory(path)
	if err != nil || isDir {
		return path
	}
	return getParentDirectory(path)
}

// MkdirIfNecessary ...
func MkdirIfNecessary(createDir string) error {
	var path string
	var err error
	//前边的判断是否是系统的分隔符
	if os.IsPathSeparator('\\') {
		path = "\\"
	} else {
		path = "/"
	}

	s := strings.Split(createDir, path)
	startIndex := 0
	dir := ""
	if s[0] == "" {
		startIndex = 1
	} else {
		dir, _ = os.Getwd() //当前的目录
	}
	for i := startIndex; i < len(s); i++ {
		d := dir + path + strings.Join(s[startIndex:i+1], path)
		if _, e := os.Stat(d); os.IsNotExist(e) {
			//在当前目录下生成md目录
			err = os.Mkdir(d, os.ModePerm)
			if err != nil {
				break
			}
		}
	}
	return err
}
