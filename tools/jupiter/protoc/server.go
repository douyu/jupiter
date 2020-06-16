package protoc

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/util/xcolor"
	template2 "github.com/douyu/jupiter/tools/jupiter/protoc/template"
	"github.com/emicklei/proto"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// GRPCServerGen  ...
type GRPCServerGen struct {
	rpcMeta *RPCMeta
}

//RPCMeta ...
type RPCMeta struct {
	Service *proto.Service
	Message []*proto.Message
	RPC     []*proto.RPC
	Package *proto.Package
	Prefix  string //生成的GRPC server代码 import pb.go文件时的前缀
}

//NewGRPCServerGen construct a GRPCServerGen instance
func NewGRPCServerGen() *GRPCServerGen {
	return &GRPCServerGen{rpcMeta: &RPCMeta{}}
}

// Parse IDL files
func (server *GRPCServerGen) parseProtoFile(protoFilePath string) (err error) {
	reader, err := os.Open(protoFilePath)
	defer reader.Close()
	if err != nil {
		return
	}
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return
	}
	proto.Walk(definition,
		proto.WithService(server.handleService),
		proto.WithMessage(server.handleMessage),
		proto.WithRPC(server.handleRPC),
		proto.WithPackage(server.handlePackage))
	return
}
func (server *GRPCServerGen) handlePackage(pkg *proto.Package) {
	server.rpcMeta.Package = pkg
}
func (server *GRPCServerGen) handleRPC(rpc *proto.RPC) {
	server.rpcMeta.RPC = append(server.rpcMeta.RPC, rpc)
}
func (server *GRPCServerGen) handleService(s *proto.Service) {
	server.rpcMeta.Service = s
}
func (server *GRPCServerGen) handleMessage(m *proto.Message) {
	server.rpcMeta.Message = append(server.rpcMeta.Message, m)
}
func (server *GRPCServerGen) generateServer() (err error) {
	if server.rpcMeta == nil {
		return
	}
	if err = server.initPrefix(); err != nil {
		return
	}
	if err = server.parseProtoFile(option.protoFilePath); err != nil {
		return
	}
	outPut := filepath.Join(option.outputDir, server.rpcMeta.Package.Name)
	if err = os.MkdirAll(outPut, 0755); err != nil {
		return
	}
	fileName := fmt.Sprintf("%sServer.go", server.rpcMeta.Service.Name)
	fileName = Lcfirst(fileName)
	filePath := path.Join(outPut, fileName)
	file, err1 := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err1 != nil {
		err = err1
		return
	}
	defer file.Close()
	filePath, _ = filepath.Abs(filePath)
	err1 = server.render(file, template2.GRPCServerTemplate, server.rpcMeta)
	if err != nil {
		err = err1
		return
	}
	fmt.Println(xcolor.Greenf("GRCP server file generate success ,the path :", filePath))
	return
}
func (server *GRPCServerGen) initPrefix() (err error) {
	if option.prefix != "" {
		server.rpcMeta.Prefix = option.prefix
		return
	}
	goPath := os.Getenv("GOPATH")
	execFilePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}
	if runtime.GOOS == "windows" {
		execFilePath = strings.Replace(execFilePath, "\\", "/", -1)
		goPath = strings.Replace(goPath, "\\", "/", -1)
	}
	lastIdx := strings.LastIndex(execFilePath, "/")
	if lastIdx < 0 {
		return
	}
	output := strings.ToLower(execFilePath[0:lastIdx])
	srcPath := path.Join(goPath, "src/")
	if srcPath[len(srcPath)-1] != '/' {
		srcPath = fmt.Sprintf("%s/", srcPath)
	}
	server.rpcMeta.Prefix = strings.Replace(output, strings.ToLower(srcPath), "", -1)
	fmt.Printf("rpcMeta.Prefix:%s\n", server.rpcMeta.Prefix)
	return
}
func (server *GRPCServerGen) render(file *os.File, data string, rpcMeta *RPCMeta) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	err = t.Execute(file, rpcMeta)
	if err != nil {
		return
	}
	return
}
