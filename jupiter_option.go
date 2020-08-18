package jupiter

import "github.com/douyu/jupiter/pkg/conf"

type Option func(a *Application)

type Disable int

const (
	DisableParserFlag      Disable = 1
	DisableLoadConfig      Disable = 2
	DisableDefaultGovernor Disable = 3
)

func (a *Application) WithOptions(options ...Option) {
	for _, option := range options {
		option(a)
	}
}

func WithConfigParser(unmarshaller conf.Unmarshaller) Option {
	return func(a *Application) {
		a.configParser = unmarshaller
	}
}

func WithDisable(d Disable) Option {
	return func(a *Application) {
		a.disableMap[d] = true
	}
}
