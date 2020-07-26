package pkg

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/util/xcolor"
)

const jupiterVersion = "0.2.0"

var (
	startTime string
	goVersion string
)

// build info
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
		appName = os.Getenv("APP_NAME")
		if appName == "" {
			appName = filepath.Base(os.Args[0])
		}
	}

	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	startTime = xtime.TS.Format(time.Now())
	goVersion = runtime.Version()
	InitEnv()
}

// Name gets application name.
func Name() string {
	return appName
}

func SetName(appName string) {
	appName = appName
}

func AppID() string {
	return appID
}

func SetAppID(appID string) {
	appID = appID
}

func AppVersion() string {
	return buildAppVersion
}

func SetAppVersion(appVersion string) {
	appVersion = appVersion
}

func JupiterVersion() string {
	return jupiterVersion
}

func SetJupiterVersion(jupiterVersion string) {
	jupiterVersion = jupiterVersion
}

func BuildTime() string {
	return buildTime
}

func BuildUser() string {
	return buildUser
}

func BuildHost() string {
	return buildHost
}

func SetBuildTime(buildTime string) {
	buildTime = strings.Replace(buildTime, "--", " ", 1)
}

// HostName ...
func HostName() string {
	return hostName
}

func StartTime() string {
	return startTime
}

func GoVersion() string {
	return goVersion
}

// PrintVersion ...
func PrintVersion() {
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("name"), xcolor.Blue(appName))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("appID"), xcolor.Blue(appID))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("appVersion"), xcolor.Blue(buildAppVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("jupiterVersion"), xcolor.Blue(jupiterVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildUser"), xcolor.Blue(buildUser))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildHost"), xcolor.Blue(buildHost))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildTime"), xcolor.Blue(BuildTime()))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildStatus"), xcolor.Blue(buildStatus))
}
