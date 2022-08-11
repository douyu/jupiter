
# JUPITER: Governance-oriented Microservice Framework

![logo](doc/logo.png)

[![GoTest](https://github.com/douyu/jupiter/workflows/Go/badge.svg)](https://github.com/douyu/jupiter/actions)
[![codecov](https://codecov.io/gh/douyu/jupiter/branch/master/graph/badge.svg)](https://codecov.io/gh/douyu/jupiter)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/douyu/jupiter?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/douyu/jupiter)](https://goreportcard.com/report/github.com/douyu/jupiter)
![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

## Introduction

JUPITER is a governance-oriented microservice framework, which is being used for years at [Douyu](https://www.douyu.com).

## Documentation

See the [中文文档](http://jupiter.douyu.com/) for the Chinese documentation.

## Quick Start

```bash
go install github.com/douyu/jupiter/cmd/jupiter@latest
jupiter new example-go
cd example-go
jupiter run -c cmd/exampleserver/.jupiter.toml
```

More Example:

- [Quick Start](doc/wiki-cn/quickstart.md)  
- [Examples](https://github.com/douyu/jupiter-examples)
- [Project Layout](https://github.com/douyu/jupiter-layout)

## Bugs and Feedback

For bug report, questions and discussions please submit an issue.

## Contributing

Contributions are always welcomed! Please see [CONTRIBUTING](CONTRIBUTING.md) for detailed guidelines.

You can start with the issues labeled with good first issue.

<a href="https://github.com/douyu/jupiter/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=douyu/jupiter" />
</a>

## Contact

- [Wechat Group](https://jupiter.douyu.com/join/#%E5%BE%AE%E4%BF%A1)
