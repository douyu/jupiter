package pkg

import (
	"crypto/md5"
	"fmt"
	"github.com/douyu/jupiter/pkg/constant"
	"os"
)

var (
	appLogDir   string
	appMode     string
	appRegion   string
	appZone     string
	appHost     string
	appInstance string
)

func InitEnv() {
	appLogDir = os.Getenv(constant.EnvAppLogDir)
	appMode = os.Getenv(constant.EnvAppMode)
	appRegion = os.Getenv(constant.EnvAppRegion)
	appZone = os.Getenv(constant.EnvAppZone)
	appHost = os.Getenv(constant.EnvAppHost)
	appInstance = os.Getenv(constant.EnvAppInstance)
	if appInstance == "" {
		appInstance = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s", HostName(), AppID()))))
	}
}

func AppLogDir() string {
	return appLogDir
}

func SetAppLogDir(appLogDir string) {
	appLogDir = appLogDir
}

func AppMode() string {
	return appMode
}

func SetAppMode(appMode string) {
	appMode = appMode
}

func AppRegion() string {
	return appRegion
}

func SetAppRegion(appRegion string) {
	appRegion = appRegion
}

func AppZone() string {
	return appZone
}

func SetAppZone(appZone string) {
	appZone = appZone
}

func AppHost() string {
	return appHost
}

func SetAppHost(appHost string) {
	appHost = appHost
}

func AppInstance() string {
	return appInstance
}
