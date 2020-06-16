package protoc

import (
	"errors"
	"fmt"
	"github.com/douyu/jupiter/pkg/util/xcolor"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	_genGoFastAddress = "go get -u -v github.com/gogo/protobuf/protoc-gen-gofast"
	_grpcProtocCmd    = "protoc --gofast_out=plugins=grpc:%s %s"
)

func generateGRPC() (err error) {
	if err = installGRPCGen(); err != nil {
		return
	}
	if err = doGenerate(); err != nil {
		return
	}
	return
}
func installGRPCGen() (err error) {
	gofastPath := ""
	if gofastPath, err = exec.LookPath("protoc-gen-gofast"); err != nil {
		fmt.Println(xcolor.Green("start installing protoc-gen-gofast"))
		if err = executeGoGet(_genGoFastAddress); err != nil {
			return
		}
	}
	fmt.Println(xcolor.Greenf("protoc-gen-gofast installation was successful, the installation path is:", gofastPath))
	return
}
func executeGoGet(address string) error {
	args := strings.Split(address, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func protocEnvCheck() (err error) {
	protocPath := ""
	if protocPath, err = exec.LookPath("protoc"); err != nil {
		err = errors.New("You haven't installed Protobuf yet，Please visit this page to install with your own system：https://github.com/protocolbuffers/protobuf/releases")
		return
	}
	fmt.Println(xcolor.Greenf("Protoc environment monitoring is successful , the installation path is:", protocPath))
	return
}
func doGenerate() (err error) {
	if err = os.MkdirAll(option.outputDir, 0755); err != nil {
		return
	}
	cmdLine := fmt.Sprintf(_grpcProtocCmd, option.outputDir, option.protoFilePath)
	fmt.Println("protocCmdLine:", cmdLine)
	args := strings.Split(cmdLine, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	if genAbsPath, err := filepath.Abs(option.outputDir); err == nil {
		fmt.Println(xcolor.Greenf("pb.go file generated successfully. The path is as follows:", genAbsPath))
	}
	return
}
