package xxl

// 通用响应
type res struct {
	Error int64       `json:"error"` // 200 表示正常、其他失败
	Msg   interface{} `json:"msg"`   // 错误提示消息
}

/*****************  上行参数  *********************/

// 注册参数
type Registry struct {
	RegistryGroup string `json:"registryGroup"`
	RegistryKey   string `json:"registryKey"`
	RegistryValue string `json:"registryValue"`
}

// 执行器执行完任务后，回调任务结果时使用
type call []*callElement

type callElement struct {
	LogID         int64          `json:"logId"`
	LogDateTim    int64          `json:"logDateTim"`
	ExecuteResult *ExecuteResult `json:"executeResult"`
}

//任务执行结果 200 表示任务执行正常，500表示失败
type ExecuteResult struct {
	Error int64       `json:"error"`
	Msg   interface{} `json:"msg"`
}

/*****************  下行参数  *********************/

// 阻塞处理策略
const (
	SerialExecution     = "SERIAL_EXECUTION"       //单机串行
	DiscardLater        = "DISCARD_LATER"          //丢弃后续调度
	DiscardLaterNoAlarm = "DISCARD_LATER_NO_ALARM" //丢弃后续调度,并不报警
	CoverEarly          = "COVER_EARLY"            //覆盖之前调度
)

// 终止任务请求参数
type killReq struct {
	JobID int64 `json:"jobId"` // 任务ID
}

// 日志请求
type logReq struct {
	LogDateTime int64 `json:"logDateTime"` // 本次调度日志时间
	LogID       int64 `json:"logId"`       // 本次调度日志ID
	FromLineNum int32 `json:"fromLineNum"` // 日志开始行号，滚动加载日志
}

// 日志响应
type logRes struct {
	Error int64     `json:"error"` // 200 表示正常、其他失败
	Msg   string    `json:"msg"`   // 错误提示消息
	Data  LogResult `json:"data"`  // 日志响应内容
}

// 日志响应内容
type LogResult struct {
	FromLineNum int32  `json:"fromLineNum"` // 本次请求，日志开始行数
	ToLineNum   int32  `json:"toLineNum"`   // 本次请求，日志结束行号
	LogContent  string `json:"logContent"`  // 本次请求日志内容
	IsEnd       bool   `json:"isEnd"`       // 日志是否全部加载完
}

// idleBeat请求
type idleReq struct {
	JobID int64 `json:"jobId"` // 任务ID
}

// 心跳响应
type heatBeatRes struct {
	Error int64       `json:"error"` // 200 表示正常、其他失败
	Msg   interface{} `json:"msg"`   // 错误提示消息
	Data  string      `json:"data"`  // 心跳数据
}

type idleData struct {
	Ip   string `json:"ip"`   // 节点ip
	Idle bool   `json:"idle"` // 是否空闲
}

type beatData struct {
	Ip string `json:"ip"` // 节点ip
}
