package jupiter

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
)

type Option func(a *Application)

type Disable int

const (
	DisableParserFlag Disable = iota + 1
	DisableLoadConfig
	DisableDefaultGovernor
	DisableBanner
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

func WithLogger(logger *xlog.Logger) Option {
	return func(a *Application) {
		a.logger = logger
	}
}

func WithDisable(d ...Disable) Option {
	return func(a *Application) {
		if len(d) == 0 {
			return
		}
		if a.disableMap == nil {
			a.disableMap = make(map[Disable]bool)
		}
		for _, disabled := range d {
			a.disableMap[disabled] = true
		}
	}
}

func WithBanner(banner string) Option {
	return func(a *Application) {
		a.banner = banner
	}
}
