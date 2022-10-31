package xxl

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/douyu/jupiter/pkg/executor"
)

func Test_Int64ToStr(t *testing.T) {
	tests := []struct {
		name string
		args int64
		want string
	}{
		{
			name: "正数",
			args: 20220320,
			want: "20220320",
		},
		{
			name: "负数",
			args: -1024,
			want: "-1024",
		},
		{
			name: "0",
			args: 0,
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64ToStr(tt.args); got != tt.want {
				t.Errorf("Int64ToStr() = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_StrToInt64(t *testing.T) {
	tests := []struct {
		name string
		args string
		want int64
	}{
		{
			name: "正数",
			args: "20220320",
			want: 20220320,
		},
		{
			name: "负数",
			args: "-1024",
			want: -1024,
		},
		{
			name: "0",
			args: "0",
			want: 0,
		},
		{
			name: "非数字",
			args: "nonum",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrToInt64(tt.args); got != tt.want {
				t.Errorf("Int64ToStr() = %d, want %d", got, tt.want)
			}
		})
	}
}

func Test_returnCall(t *testing.T) {
	type Args struct {
		param *executor.RunReq
		code  int64
		msg   string
	}
	tests := []struct {
		name string
		args Args
		want []byte
	}{
		{
			name: "returnCall单元测试",
			args: Args{
				param: &executor.RunReq{
					JobID:       123,
					LogID:       101,
					LogDateTime: 123456,
				},
				code: http.StatusOK,
				msg:  "There are tasks running",
			},
			want: []byte("[{\"logId\":101,\"logDateTim\":123456,\"executeResult\":{\"error\":200,\"msg\":\"There are tasks running\"}}]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnCall(tt.args.param, tt.args.code, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnCall() = %d, want %d", got, tt.want)
			}
		})
	}
}

func Test_returnKill(t *testing.T) {
	type Args struct {
		param *killReq
		code  int64
	}
	tests := []struct {
		name string
		args Args
		want []byte
	}{
		{
			name: "returnKill单元测试",
			args: Args{
				param: &killReq{
					JobID: 123,
				},
				code: http.StatusInternalServerError,
			},
			want: []byte("{\"error\":500,\"msg\":\"Task kill err\"}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnKill(tt.args.param, tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnKill() = %s, want %d", string(got), tt.want)
			}
		})
	}
}

func Test_returnLog(t *testing.T) {
	type Args struct {
		param *logReq
		code  int64
	}
	tests := []struct {
		name string
		args Args
		want []byte
	}{
		{
			name: "returnLog单元测试",
			args: Args{
				param: &logReq{
					LogID:       123,
					FromLineNum: 10,
				},
				code: http.StatusOK,
			},
			want: []byte("{\"error\":200,\"msg\":\"success\",\"data\":{\"fromLineNum\":10,\"toLineNum\":1,\"logContent\":\"\",\"isEnd\":true}}"),
		},
		{
			name: "returnLog单元测试",
			args: Args{
				param: &logReq{
					LogID:       123,
					FromLineNum: 10,
				},
				code: http.StatusBadRequest,
			},
			want: []byte("{\"error\":400,\"msg\":\"log err\",\"data\":{\"fromLineNum\":10,\"toLineNum\":1,\"logContent\":\"\",\"isEnd\":true}}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnLog(tt.args.param, tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnLog() = %s, want %d", string(got), tt.want)
			}
		})
	}
}

func Test_returnGeneral(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "returnGeneral单元测试",
			want: []byte("{\"error\":200,\"msg\":\"\"}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnGeneral(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnGeneral() = %s, want %d", string(got), tt.want)
			}
		})
	}
}

func Test_returnHeatBeat(t *testing.T) {
	str, _ := json.Marshal(beatData{
		Ip: "127.0.0.1",
	})
	data := &heatBeatRes{
		Error: http.StatusOK,
		Msg:   "",
		Data:  string(str),
	}
	body, _ := json.Marshal(data)
	tests := []struct {
		name string
		args string
		want []byte
	}{
		{
			name: "returnHeatBeat单元测试",
			args: "127.0.0.1",
			want: body,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnHeatBeat(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnHeatBeat() = %s, want %d", string(got), tt.want)
			}
		})
	}
}

func Test_returnIdle(t *testing.T) {
	type Args struct {
		ip   string
		idel bool
	}
	tests := []struct {
		name string
		args Args
		want []byte
	}{
		{
			name: "returnIdle单元测试",
			args: Args{
				ip:   "127.0.0.1",
				idel: true,
			},
			want: []byte("{\"error\":200,\"msg\":\"\",\"data\":\"{\\\"ip\\\":\\\"127.0.0.1\\\",\\\"idle\\\":true}\"}"),
		},
		{
			name: "returnIdle单元测试",
			args: Args{
				ip:   "127.0.0.1",
				idel: false,
			},
			want: []byte("{\"error\":500,\"msg\":\"\",\"data\":\"\"}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnIdle(tt.args.ip, tt.args.idel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnIdle() = %s, want %d", string(got), tt.want)
			}
		})
	}
}

func Test_returnAuthError(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "returnAuthError单元测试",
			want: []byte("{\"error\":500,\"msg\":\"auth token error\"}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnAuthError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnAuthError() = %s, want %d", string(got), tt.want)
			}
		})
	}
}
