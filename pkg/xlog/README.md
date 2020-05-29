# xlog

xlog wrapped go.uber.org/zap, simplify the difficulty of use.


## 动态设置日志级别

修改默认日志级别:
```toml
[jupiter.logger.default]
    level = "error"
```
修改自定义日志界别:
```toml
[jupiter.logger.mylog]
    level = "error"
```

## 创建自定义日志

```golang
logger := xlog.StdConfig("mylog").Build()
logger.Info("info", xlog.String("a", "b"))
logger.Infof("info %s", "a")
logger.Infow("info", "a", "b")
```

也可以更精确的控制:
```golang
config := xlog.Config{
    Name: "default.log",
    Dir: "/tmp",
    Level: "info",
}
logger := config.Build()
logger.SetLevel(xlog.DebugLevel)
logger.Debug("debug", xlog.String("a", "b"))
logger.Debugf("debug %s", "a")
logger.Debugw("debug", "a", "b")
```

