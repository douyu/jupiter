module github.com/douyu/jupiter

go 1.24.0

require (
	cirello.io/pglock v1.16.1
	github.com/BurntSushi/toml v1.6.0
	github.com/alibaba/sentinel-golang v1.0.4
	github.com/aliyun/aliyun-tablestore-go-sdk v1.7.17
	github.com/apache/rocketmq-client-go/v2 v2.1.3-0.20250427084711-67ec50b93040
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0
	github.com/coocood/freecache v1.2.4
	github.com/cosmtrek/air v1.45.0
	github.com/davecgh/go-spew v1.1.1
	github.com/dimiro1/banner v1.1.0
	github.com/fatih/color v1.18.0
	github.com/fsnotify/fsnotify v1.9.0
	github.com/gin-gonic/gin v1.11.0
	github.com/go-basic/ipv4 v1.0.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-resty/resty/v2 v2.17.1
	github.com/gogf/gf v1.16.9
	github.com/golang/protobuf v1.5.4
	github.com/gorilla/websocket v1.5.3
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.7
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/hnlq715/struct2interface v0.1.5
	github.com/iancoleman/strcase v0.3.0
	github.com/json-iterator/go v1.1.12
	github.com/juju/ratelimit v1.0.2
	github.com/labstack/echo/v4 v4.15.0
	github.com/mattn/go-colorable v0.1.14
	github.com/mitchellh/mapstructure v1.5.0
	github.com/modern-go/reflect2 v1.0.2
	github.com/onsi/ginkgo/v2 v2.28.1
	github.com/onsi/gomega v1.39.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/philchia/agollo/v4 v4.1.5
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.20.5
	github.com/robfig/cron/v3 v3.0.1
	github.com/samber/lo v1.52.0
	github.com/smallnest/weighted v0.0.0-20200122032019-adf21c9b8bd1
	github.com/smartystreets/goconvey v1.8.1
	github.com/spf13/cast v1.10.0
	github.com/srikrsna/protoc-gen-gotag v1.0.2
	github.com/stretchr/testify v1.11.1
	github.com/tidwall/pretty v1.2.1
	github.com/urfave/cli v1.22.17
	github.com/valyala/fasthttp v1.69.0
	github.com/xlab/treeprint v1.2.0
	go.etcd.io/etcd/api/v3 v3.6.7
	go.etcd.io/etcd/client/v3 v3.6.7
	go.mongodb.org/mongo-driver v1.17.8
	go.opentelemetry.io/otel v1.39.0
	go.opentelemetry.io/otel/bridge/opentracing v1.39.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.39.0
	go.opentelemetry.io/otel/sdk v1.39.0
	go.opentelemetry.io/otel/trace v1.39.0
	go.uber.org/automaxprocs v1.6.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/mod v0.32.0
	golang.org/x/sync v0.19.0
	golang.org/x/text v0.33.0
	google.golang.org/genproto/googleapis/api v0.0.0-20260128011058-8636f8732409
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260128011058-8636f8732409
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.30.0
)

exclude github.com/aliyun/aliyun-tablestore-go-sdk v4.1.3+incompatible

require (
	dario.cat/mergo v1.0.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bep/godartsass v0.16.0 // indirect
	github.com/bep/golibsass v1.1.0 // indirect
	github.com/bytedance/sonic v1.14.0 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clbanning/mxj v1.8.5-0.20200714211355-ff02cfb8ea28 // indirect
	github.com/cli/safeexec v1.0.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/creack/pty v1.1.18 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gohugoio/hugo v0.111.3 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20260115054156-294ebfa9ad83 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/grokify/html-strip-tags-go v0.0.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil/v3 v3.21.7 // indirect
	github.com/sirupsen/logrus v1.8.3 // indirect
	github.com/smarty/assertions v1.15.0 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/tdewolff/parse/v2 v2.6.5 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tklauser/go-sysconf v0.3.7 // indirect
	github.com/tklauser/numcpus v0.2.3 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.6.7 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/mock v0.5.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/arch v0.20.0 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/exp v0.0.0-20221031165847-c99f073a8326 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	stathat.com/c/consistent v1.0.0 // indirect
)
