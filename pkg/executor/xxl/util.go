package xxl

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/douyu/jupiter/pkg/executor"
	"github.com/douyu/jupiter/pkg/executor/xxl/logger"
)

//int64 to str
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

//str to int64
func StrToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

//执行任务回调
func returnCall(req *executor.RunReq, code int64, msg string) []byte {
	data := call{
		&callElement{
			LogID:      req.LogID,
			LogDateTim: req.LogDateTime,
			ExecuteResult: &ExecuteResult{
				Error: code,
				Msg:   msg,
			},
		},
	}
	str, _ := json.Marshal(data)
	return str
}

//杀死任务返回
func returnKill(req *killReq, code int64) []byte {
	msg := ""
	if code != http.StatusOK {
		msg = "Task kill err"
	}
	data := res{
		Error: code,
		Msg:   msg,
	}
	str, _ := json.Marshal(data)
	return str
}

//日志返回
func returnLog(req *logReq, code int64) []byte {
	msg := "success"
	if code != http.StatusOK {
		msg = "log err"
	}
	line, content := logger.ReadLog(req.LogDateTime, req.LogID, req.FromLineNum)
	logResult := LogResult{
		FromLineNum: req.FromLineNum,
		ToLineNum:   line,
		LogContent:  content,
		IsEnd:       true,
	}
	data := &logRes{Error: code, Msg: msg, Data: logResult}
	str, _ := json.Marshal(data)
	return str
}

//通用返回
func returnGeneral() []byte {
	data := &res{
		Error: http.StatusOK,
		Msg:   "",
	}
	str, _ := json.Marshal(data)
	return str
}

//心跳返回
func returnHeatBeat(ip string) []byte {
	str, _ := json.Marshal(beatData{
		Ip: ip,
	})
	data := &heatBeatRes{
		Error: http.StatusOK,
		Msg:   "",
		Data:  string(str),
	}
	body, _ := json.Marshal(data)
	return body
}

//idle返回
func returnIdle(ip string, idle bool) []byte {
	if !idle {
		data := &heatBeatRes{
			Error: http.StatusInternalServerError,
			Msg:   "",
		}
		body, _ := json.Marshal(data)
		return body
	}
	str, _ := json.Marshal(idleData{
		Ip:   ip,
		Idle: idle,
	})
	data := &heatBeatRes{
		Error: http.StatusOK,
		Msg:   "",
		Data:  string(str),
	}
	body, _ := json.Marshal(data)
	return body
}

//通用返回
func returnAuthError() []byte {
	data := &res{
		Error: http.StatusInternalServerError,
		Msg:   "auth token error",
	}
	str, _ := json.Marshal(data)
	return str
}
