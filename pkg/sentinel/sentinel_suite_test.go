package sentinel_test

import (
	"testing"

	"github.com/alibaba/sentinel-golang/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSentinel(t *testing.T) {
	logging.ResetGlobalLoggerLevel(logging.DebugLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sentinel Suite")
}
