### uuid server 
## 特性：
- 支持uuid 的生成
- 全局唯一
- 可通过配置文件，flag 配置
- 可通过本地手动配置唯一 nodeId
- 可通过redis 配置全局唯一 NodeId

## 配置文件
通过配置这些字段可以控制 uuid 的组成部分
```toml
[jupiter.server.uuid]
    epoch = 1288834974657
    nodeBits = 10
    stepBits = 12
    nodeId = 1
    enableRedis = true  # 通过redis 来配置NodeId，配置文件的NodeId将无效
```

通过这这属性来配置redis的地址
```toml
[jupiter.redis.uuid]
    addrs = ["10.1.0.86:6379"]
```
