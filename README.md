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

## Requirements

- Go version >= 1.19
- Docker

## Quick Start

1. Install [jupiter](https://github.com/douyu/jupiter/tree/master/cmd/jupiter) toolkit
1. Create example project from [jupiter-layout](https://github.com/douyu/jupiter-layout)
1. Download go mod dependencies
1. Run the example project with [jupiter](https://github.com/douyu/jupiter/tree/master/cmd/jupiter) toolkit
1. Just code yourself :-)

```bash
go install github.com/douyu/jupiter/cmd/jupiter@latest
jupiter new example-go
cd example-go
go mod tidy
docker-compose -f tests/e2e/docker-compose.yml up -d
jupiter run -c cmd/exampleserver/.jupiter.toml
```

More Example:

- [Project Layout](https://github.com/douyu/jupiter-layout)
- [Examples](https://github.com/douyu/jupiter-examples)

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
