package pkg

import (
	"crypto/md5"
	"fmt"
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

var (
	EnvAppLogDir   = "APP_LOG_DIR"
	EnvAppMode     = "APP_MODE"
	EnvAppRegion   = "APP_REGION"
	EnvAppZone     = "APP_ZONE"
	EnvAppHost     = "APP_HOST"
	EnvAppInstance = "APP_INSTANCE" // application unique instance id.
)

func InitEnv() {
	appLogDir = os.Getenv(EnvAppLogDir)
	appMode = os.Getenv(EnvAppMode)
	appRegion = os.Getenv(EnvAppRegion)
	appZone = os.Getenv(EnvAppZone)
	appHost = os.Getenv(EnvAppHost)
	appInstance = os.Getenv(EnvAppInstance)
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
