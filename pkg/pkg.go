// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/util/xcolor"
)

const jupiterVersion = "0.2.0"

var (
	startTime string
	goVersion string
)

// build info
/*

 */
var (
	appName         string
	appID           string
	hostName        string
	buildAppVersion string
	buildUser       string
	buildHost       string
	buildStatus     string
	buildTime       string
)

func init() {
	if appName == "" {
		appName = os.Getenv(constant.EnvAppName)
		if appName == "" {
			appName = filepath.Base(os.Args[0])
		}
	}

	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	startTime = time.Now().Format("2006-01-02 15:04:05")
	SetBuildTime(buildTime)
	goVersion = runtime.Version()
	InitEnv()
}

// Name gets application name.
func Name() string {
	return appName
}

//SetName set app anme
func SetName(s string) {
	appName = s
}

//AppID get appID
func AppID() string {
	if appID == "" {
		return "1234567890" //default appid when APP_ID Env var not set
	}
	return appID
}

//SetAppID set appID
func SetAppID(s string) {
	appID = s
}

//AppVersion get buildAppVersion
func AppVersion() string {
	return buildAppVersion
}

// SetAppVersion set appVersion
func SetAppVersion(s string) {
	buildAppVersion = s
}

//JupiterVersion get jupiterVersion
func JupiterVersion() string {
	return jupiterVersion
}

// todo: jupiterVersion is const not be set
// func SetJupiterVersion(s string) {
// 	jupiterVersion = s
// }

//BuildTime get buildTime
func BuildTime() string {
	return buildTime
}

//BuildUser get buildUser
func BuildUser() string {
	return buildUser
}

//BuildHost get buildHost
func BuildHost() string {
	return buildHost
}

//SetBuildTime set buildTime
func SetBuildTime(param string) {
	buildTime = strings.Replace(param, "--", " ", 1)
}

// HostName get host name
func HostName() string {
	return hostName
}

//StartTime get start time
func StartTime() string {
	return startTime
}

//GoVersion get go version
func GoVersion() string {
	return goVersion
}

func LogDir() string {
	// LogDir gets application log directory.
	logDir := AppLogDir()
	if logDir == "" {
		if appPodIP != "" && appPodName != "" {
			// k8s 环境
			return fmt.Sprintf("/home/www/logs/applogs/%s/%s/", Name(), appPodName)
		}
		return fmt.Sprintf("/home/www/logs/applogs/%s/%s/", Name(), appInstance)
	}
	return fmt.Sprintf("%s/%s/%s/", logDir, Name(), appInstance)
}

// PrintVersion print formated version info
func PrintVersion() {
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("name"), xcolor.Blue(appName))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("appID"), xcolor.Blue(appID))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("region"), xcolor.Blue(AppRegion()))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("zone"), xcolor.Blue(AppZone()))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("appVersion"), xcolor.Blue(buildAppVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("jupiterVersion"), xcolor.Blue(jupiterVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildUser"), xcolor.Blue(buildUser))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildHost"), xcolor.Blue(buildHost))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildTime"), xcolor.Blue(BuildTime()))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildStatus"), xcolor.Blue(buildStatus))
}
