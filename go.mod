module github.com/douyu/jupiter

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/rocketmq-client-go/v2 v2.0.0
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dustin/go-humanize v0.0.0-20171111073723-bb3d318650d4 // indirect
	github.com/emicklei/proto v1.9.0 // indirect
	github.com/fatih/structtag v1.2.0
	github.com/flosch/pongo2 v0.0.0-20200518135938-dfb43dbdc22a
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-resty/resty/v2 v2.2.0
	github.com/gobuffalo/packr/v2 v2.2.0 // indirect
	github.com/gogf/gf v1.13.3
	github.com/golang/mock v1.4.4 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.9.6 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.10
	github.com/labstack/echo/v4 v4.1.16
	github.com/mitchellh/mapstructure v1.3.2
	github.com/modern-go/reflect2 v1.0.1
	github.com/onsi/ginkgo v1.12.3 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/philchia/agollo/v4 v4.0.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.6.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/smallnest/weighted v0.0.0-20200122032019-adf21c9b8bd1
	github.com/smartystreets/assertions v0.0.0-20190401211740-f487f9de1cd3 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.6.1
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	github.com/tidwall/pretty v1.0.1
	github.com/uber/jaeger-client-go v2.23.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/urfave/cli v1.22.5
	go.mongodb.org/mongo-driver v1.4.1
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/multierr v1.5.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392 // indirect
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20201126233918-771906719818 // indirect
	golang.org/x/tools v0.0.0-20200806022845-90696ccdc692 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20200122232147-0452cf42e150
	google.golang.org/grpc v1.26.0
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
