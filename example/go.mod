module github.com/douyu/jupiter/example

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/rocketmq-client-go/v2 v2.0.0
	github.com/douyu/jupiter v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.7.0
	github.com/gogf/gf v1.13.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/websocket v1.4.2
	github.com/labstack/echo/v4 v4.1.16
	github.com/sentinel-group/sentinel-go-adapters v1.0.1
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	go.etcd.io/etcd/client/v3 v3.5.0
	go.mongodb.org/mongo-driver v1.5.1
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)

replace (
	github.com/douyu/jupiter => ../
	google.golang.org/grpc => google.golang.org/grpc v1.40.0
)
