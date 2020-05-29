# conf

### 从字符串中加载配置

```golang
var content = `[app] mode="dev"`
if err := conf.Load(bytes.NewBufferString(content), toml.Unmarshal); err != nil {
    panic(err)
}
```

### 从配置文件中加载配置

```golang
import (
    file_datasource "github.com/douyu/jupiter/pkg/datasource/file"
)

provider := file_datasource.NewDataSource(path)
if err := conf.Load(provider, toml.Unmarshal); err != nil {
    panic(err)
}
```

### 从etcd中加载配置

```golang
import (
   etcdv3_datasource "github.com/douyu/jupiter/pkg/datasource/etcdv3"
   "github.com/douyu/jupiter/client/etcdv3"
)
provider := etcdv3_datasource.NewDataSource(
    etcdv3.StdConfig("config_datasource").Build(),
    "/config/my-application",
)
if err := conf.Load(provider, json.Unmarshal); err != nil {
    panic(err)
}
```