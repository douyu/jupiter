package xxl

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/executor"
)

// 测试XJob
type TestXJob2 struct{}

func (ins *TestXJob2) GetJobName() string {
	return "test"
}

func (ins *TestXJob2) Run(ctx context.Context, param *executor.RunReq) (msg string, err error) {
	fmt.Println("test2")
	return "success", nil
}

var callback = func(ctx context.Context, status int, msg string) error {
	return nil
}

//TODO: Implement more test
func Test_Run2(t *testing.T) {
	job := &TestXJob2{}
	ctx, cancel := context.WithCancel(context.Background())
	task := &Task{
		Id:   1,
		Name: job.GetJobName(),
		fn:   job.Run,
		Param: &executor.RunReq{
			JobID:          123,
			ExecutorParams: "",
		},
		ctx:    ctx,
		cancel: cancel,
	}
	tests := []struct {
		name string
		want string
	}{
		{
			name: "TaskRun",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task.Run(ctx, callback)
		})
	}
}

func Test_Cancel(t *testing.T) {
	job := &TestXJob2{}
	ctx, cancel := context.WithCancel(context.Background())
	task := &Task{
		Id:   1,
		Name: job.GetJobName(),
		fn:   job.Run,
		Param: &executor.RunReq{
			JobID:          123,
			ExecutorParams: "",
		},
		cancel: cancel,
		ctx:    ctx,
	}
	task2 := &Task{
		Id:   1,
		Name: job.GetJobName(),
		fn:   job.Run,
		Param: &executor.RunReq{
			JobID:          123,
			ExecutorParams: "",
		},
	}
	tests := []struct {
		name string
		args *Task
	}{
		{
			name: "Cancel1",
			args: task,
		},
		{
			name: "Cancel2",
			args: task2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.Cancel()
		})
	}
}

func Test_SetTimeout(t *testing.T) {
	task := &Task{
		Id: 1,
	}
	tests := []struct {
		name string
		args time.Duration
		want *Task
	}{
		{
			name: "SetTimeout",
			args: 3 * time.Second,
			want: &Task{
				Id:      1,
				timeout: 3 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if task.SetTimeout(tt.args); !reflect.DeepEqual(task, tt.want) {
				t.Errorf("SetTimeout() failed")
			}
		})
	}
}

func Test_IsRunning(t *testing.T) {
	task := &Task{
		Id: 1,
	}
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "IsRunning",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := task.IsRunning(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsRunning() failed")
			}
		})
	}
}

func Test_GetParam(t *testing.T) {
	task := &Task{
		Id:    1,
		Param: &executor.RunReq{LogID: 1},
	}
	tests := []struct {
		name string
		want *executor.RunReq
	}{
		{
			name: "GetParam",
			want: &executor.RunReq{LogID: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := task.GetParam(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParam() failed")
			}
		})
	}
}

func Test_GetId(t *testing.T) {
	task := &Task{
		Id:    1,
		Param: &executor.RunReq{LogID: 1},
	}
	tests := []struct {
		name string
		want int64
	}{
		{
			name: "GetId",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := task.GetId(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetId() failed")
			}
		})
	}
}

func Test_GetName(t *testing.T) {
	task := &Task{
		Id:   1,
		Name: "testtask",
	}
	tests := []struct {
		name string
		want string
	}{
		{
			name: "GetName",
			want: "testtask",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := task.GetName(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetName() failed")
			}
		})
	}
}

func Test_GetContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	task := &Task{
		Id:     1,
		cancel: cancel,
		ctx:    ctx,
	}
	tests := []struct {
		name string
		want context.Context
	}{
		{
			name: "GetContext",
			want: ctx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := task.GetContext(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContext() failed")
			}
		})
	}
}

func Test_Info(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	task := &Task{
		Id:     1,
		Name:   "test_task",
		Param:  &executor.RunReq{LogID: 1},
		cancel: cancel,
		ctx:    ctx,
	}
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Info",
			want: "任务ID[1]任务名称[test_task]参数：",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := task.Info(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Info() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Trace(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	task := &Task{
		Id:     1,
		Name:   "test_task",
		Param:  &executor.RunReq{LogID: 1},
		cancel: cancel,
		ctx:    ctx,
	}

	t.Run("trace1", func(t *testing.T) {
		task.Trace("succuss")
	})
	task = nil
	t.Run("trace1", func(t *testing.T) {
		task.Trace("fail")
	})
}
