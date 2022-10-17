# jupiter

## 脚手架介绍

### jupiter脚手架工具集

1. 快速生成模板项目
2. 快速启动并监听文件变更

# go version

 GO >= 1.18

## 脚手架工具获取

windows :  

```shell script
set GO111MODULE=on
set GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

linux :  

```shell script
export GO111MODULE=on
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

`go install -v github.com/douyu/jupiter/cmd/jupiter@latest`

* windows 用户:  
  会在${GOPATH}/src/bin 目录下生成jupiter.exe 文件,若想方便的在任何地方使用jupiter命令,请将该
  命令配置在系统的环境变量中

## 如何使用

```
$ jupiter -h
                                     
        (_)_   _ _ __ (_) |_ ___ _ __
        | | | | | '_ \| | __/ _ \ '__|
        | | |_| | |_) | | ||  __/ |
   _/ |\__,_| .__/|_|\__\___|_|
  |__/      |_|   
                                                           
GoVersion: go1.18
GOOS: linux
GOARCH: amd64
NumCPU: 12
GOPATH: /home/liqi/.gvm/pkgsets/go1.18/global
GOROOT: /home/liqi/.gvm/gos/go1.18
Compiler: gc
ENV: /home/liqi/.gvm/pkgsets/go1.18/global
Now: Friday, 13 May 2022

NAME:
   jupiter - Fast bootstrap tool for jupiter framework

USAGE:
   jupiter [global options] command [command options] [arguments...]

COMMANDS:
   clean                clear all cached
   new, n               generate code framework
   run, r               auto restart program when files changed
   update, upgrade      Upgrade to the latest version
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```
