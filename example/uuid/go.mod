module uuid

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/douyu/jupiter v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.1.2
	github.com/google/wire v0.5.0
	github.com/labstack/echo/v4 v4.6.3
	github.com/onsi/ginkgo/v2 v2.1.4
	github.com/onsi/gomega v1.19.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/douyu/jupiter => ../../