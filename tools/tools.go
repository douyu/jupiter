//go:build tools
// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/douyu/jupiter/cmd/jupiter"
	_ "github.com/douyu/jupiter/cmd/protoc-gen-go-echo"
	_ "github.com/douyu/jupiter/cmd/protoc-gen-go-gin"
	_ "github.com/douyu/jupiter/cmd/protoc-gen-go-xerror"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/onsi/ginkgo/v2/ginkgo"
	_ "github.com/srikrsna/protoc-gen-gotag"
	_ "github.com/vektra/mockery/v2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "gorm.io/gen/tools/gentool"
)
