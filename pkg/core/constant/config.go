package constant

import "strings"

var (
	// configPrefix 配置前缀
	configPrefix = "jupiter"
	// DefaultNamespace 默认命名空间
	DefaultNamespace = ""
)

func ConfigKey(key ...string) string {
	if configPrefix == "" {
		return strings.Join(key, ".")
	}

	return configPrefix + "." + strings.Join(key, ".")
}

func SetConfigPrefix(prefix string) {
	configPrefix = prefix
}

func GetConfigPrefix() string {
	return configPrefix
}
