package e2e

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/conf/datasource/file"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2ESuites(t *testing.T) {
	conf.LoadFromDataSource(file.NewDataSource("../../config/uuidserver-local-live.toml", false), toml.Unmarshal)

	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e test cases")
}
