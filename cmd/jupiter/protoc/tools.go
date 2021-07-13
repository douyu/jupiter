package protoc

import (
	"errors"
	"github.com/urfave/cli"
)

//Run ...
func Run(cli *cli.Context) (err error) {
	// 校验OS是否安装 protoc 工具
	if err = protocEnvCheck(); err != nil {
		return
	}
	if option.protoFilePath == "" {
		err = errors.New("Please specify the proto file path and use jupiter protoc - h to view the detailed prompt")
		return
	}
	if option.outputDir == "" {
		err = errors.New("Please specify the code generation path and use jupiter protoc - h to view the detailed prompt")
		return
	}
	// 默认生成全部
	if !option.withGRPC && !option.withServer {
		option.withGRPC = true
		option.withServer = true
	}
	// 根据指定目录下的proto 文件 生成pb.go 文件
	if option.withGRPC {
		if err = generateGRPC(); err != nil {
			return
		}
	}
	// 生成 grpc server 服务端代码
	if option.withServer {
		serverGen := NewGRPCServerGen()
		if err = serverGen.generateServer(); err != nil {
			return
		}
	}
	return
}
