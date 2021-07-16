package jupiter

import (
	"github.com/douyu/jupiter/pkg/application"
)

type Option = application.Option

type Disable = application.Disable

const (
	DisableParserFlag      Disable = application.DisableParserFlag
	DisableLoadConfig      Disable = application.DisableLoadConfig
	DisableDefaultGovernor Disable = application.DisableDefaultGovernor
)

var WithConfigParser = application.WithConfigParser
var WithDisable = application.WithDisable
