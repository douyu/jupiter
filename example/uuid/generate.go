package main

// 基于go generate自动生成相关代码
// 自动生成interface：https://github.com/hnlq715/struct2interface
// proto代码生成：https://github.com/bufbuild/buf
// 依赖注入：https://github.com/google/wire
// mock代码生成：https://github.com/vektra/mockery

//go:generate struct2interface -d internal/pkg
//go:generate buf generate
//go:generate wire ./...
//go:generate mockery --all --keeptree --dir internal/pkg --output gen/mocks
