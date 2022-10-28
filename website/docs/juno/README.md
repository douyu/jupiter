# Juno

> A Distributed Application Management System

Juno 译名朱诺。这个名字来源于古罗马神话中的众神之母。它是斗鱼的微服务管理系统，
如同众神之母一样守护着所有微服务系统。

## Overview

Juno 提供了对微服务进行管理的多项能力，包括应用监控、依赖分析、配置管理等。

## 定义

|  名词 | 含义 |
|:--------------|:-----|
|`节点`|  单台ECS服务器、容器里的某个pod |
|`应用`| 某个服务 |
|`可用区`| 地区-机房-环境，这三个因素组成例如：ali(阿里云-华北 2)-ALIYUN-HB2-G(阿里云-北京可用区G）-gray|


## 功能

- 工作台：平台整体概览，默认权限即可查看
- 应用服务
  - 详情，应用和节点的关联关系
  - 监控，服务监控
  - 配置，服务配置下发
  - 日志，服务日志
  - pprof，线上实时分析
  - etcd， k-v查询
  - 事件，服务事件流
- 功能看板
  - 大盘，更全面的平台数据概览
  - Grafana，内嵌grafana
  - 依赖拓扑，根据服务配置进行的依赖解析
  - 版本管理，继续go.mod文件内容获取各个服务依赖的数据版本
  - 注册信息，ETCD数据查询
- 任务中心，分布式定时任务
- 资源中心，管理应用、节点、可用区的原始数据以及相关依赖
- 配置中心，配置模板和全局资源配置
- 测试中心，提供GRPC和HTTP测试工具
- 权限管理，接入Casbin权限管理
- 系统设置

## 快速开始

### Requirements

1. Etcd
2. MySQL

### 二进制安装包和安装

```bash
## 下载
wget https://github.com/douyu/juno/release/juno.tar.gz

## 解压
tar -zxvf juno.tar.gz
```

### 初始化和启动

[参考配置](https://github.com/douyu/juno/blob/master/config/single-region-admin.toml)

```bash
cd juno

# 安装数据库
./bin/install --config=./参考配置

# 启动juno后台
./bin/juno --config=./参考配置

# 启动juno的agent
./bin/juno-agent --config=./参考配置

```

完成以上步骤后，Juno 后台会默认使用 50000 端口提供后台管理界面服务。
在浏览器中打开 `http://localhost:50000`，可以看到登录界面。

初始的用户名： `admin`

初始密码： `admin`

登录后可以看到如下界面代表安装成功：
