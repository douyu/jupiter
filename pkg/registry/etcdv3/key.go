package etcdv3

import (
	"fmt"
	"strings"

	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/util/xnet"
	"github.com/douyu/jupiter/pkg/util/xstring"
)

// RegisterKey ...
type RegisterKey struct {
	Prefix  string
	AppName string
	Kind    registry.Kind
	Scheme  string
	Host    string
}

// String ...
func (rk RegisterKey) String() string {
	return fmt.Sprintf("/%s/%s/%s/%s://%s", rk.Prefix, rk.AppName, rk.Kind, rk.Scheme, rk.Host)
}

// ToRegistryKey ...
func ToRegistryKey(key string) (*RegisterKey, error) {
	rk := &RegisterKey{}
	key = strings.TrimLeft(key, "/")
	prefix, appName, kind, addr := xstring.Split(key, "/").Head4()
	rk.Prefix = prefix
	rk.AppName = appName
	rk.Kind = registry.ToKind(kind)

	url, err := xnet.ParseURL(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid addr: %s,  %w", addr, err)
	}

	rk.Scheme = url.Scheme
	rk.Host = url.Host
	return rk, nil
}
