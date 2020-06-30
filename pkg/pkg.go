package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/util/xcolor"
)

var (
	appName          string
	hostName         string
	region           string
	zone             string
	buildVersion     string
	buildGitRevision string
	buildUser        string
	buildHost        string
	buildStatus      string
	buildTime        string
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

	initMetadata(&region, constant.EnvRegion, constant.DefaultRegion)
	initMetadata(&zone, constant.EnvZone, constant.DefaultZone)
}

// Name gets application name.
func Name() string {
	return appName
}

// HostName ...
func HostName() string {
	return hostName
}

// Region ...
func Region() string {
	return region
}

// Zone ...
func Zone() string {
	return zone
}

// PrintVersion ...
func PrintVersion() {
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("name"), xcolor.Blue(appName))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("region"), xcolor.Blue(region))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("zone"), xcolor.Blue(zone))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("version"), xcolor.Blue(buildVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("revision"), xcolor.Blue(buildGitRevision))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("user"), xcolor.Blue(buildUser))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("host"), xcolor.Blue(buildHost))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildTime"), xcolor.Blue(buildTime))
	fmt.Printf("%-8s]> %-30s => %s\n", "jupiter", xcolor.Red("buildStatus"), xcolor.Blue(buildStatus))
}

func initMetadata(val *string, envKey string, defaultValue string) {
	if *val == "" {
		*val = os.Getenv(envKey)
	}

	if *val == "" {
		*val = defaultValue
	}
}
