package xxl

import (
	"reflect"
	"testing"
)

func Test_Set(t *testing.T) {
	taskList := &taskList{
		data: make(map[string]*TaskWithPending),
	}
	task1 := &Task{
		Id: 1,
	}
	type Args struct {
		key string
		val *Task
	}
	tests := []struct {
		name string
		args Args
		want *Task
	}{
		{
			name: "set",
			args: Args{
				key: "task1",
				val: task1,
			},
			want: task1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if taskList.Set(tt.args.key, tt.args.val); !reflect.DeepEqual(taskList.Get("task1"), task1) {
				t.Errorf("Executor_List_Set() failed")
			}
		})
	}
}

func Test_Get(t *testing.T) {
	taskList := &taskList{
		data: make(map[string]*TaskWithPending),
	}
	task1 := &Task{
		Id: 1,
	}
	taskList.Set("task1", task1)
	tests := []struct {
		name string
		args string
		want *Task
	}{
		{
			name: "get",
			args: "task1",
			want: task1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := taskList.Get(tt.args); !reflect.DeepEqual(got, task1) {
				t.Errorf("Executor_List_Get() failed")
			}
		})
	}
}

//任务排入队列
func Test_EnqueuePending(t *testing.T) {
	taskList := &taskList{
		data: make(map[string]*TaskWithPending),
	}
	task1 := &Task{
		Id:      1,
		running: 0,
		Name:    "task1",
	}
	taskList.Set("task1", task1)
	t.Run("nil", func(t *testing.T) {
		if err := taskList.EnqueuePending("task1"); err != nil {
			t.Errorf("Executor_List_EnqueuePending() failed")
		}
	})
	t.Run("invalid key", func(t *testing.T) {
		if err := taskList.EnqueuePending("invalid"); err.Error() != "invalid key" {
			t.Errorf("Executor_List_EnqueuePending() failed")
		}
	})
	t.Run("over limit", func(t *testing.T) {
		_ = taskList.EnqueuePending("task1")
		if err := taskList.EnqueuePending("task1"); err.Error() != "over limit" {
			t.Errorf("Executor_List_EnqueuePending() failed")
		}
	})
}

//获取队列
func Test_GetPending(t *testing.T) {
	taskList := &taskList{
		data: make(map[string]*TaskWithPending),
	}
	task1 := &Task{
		Id:      1,
		running: 0,
		Name:    "task1",
	}
	taskList.Set("task1", task1)
	t.Run("task1", func(t *testing.T) {
		if task, _ := taskList.GetPending("task1"); !reflect.DeepEqual(task, task1) {
			t.Errorf("Test_GetPending() failed")
		}
	})
	t.Run("nil", func(t *testing.T) {
		if task, _ := taskList.GetPending("invalid"); task != nil {
			t.Errorf("Test_GetPending() failed")
		}
	})
}

func Test_GetAllPending(t *testing.T) {
	taskList := &taskList{
		data: make(map[string]*TaskWithPending),
	}
	task1 := &Task{
		Id:      1,
		running: 0,
		Name:    "task1",
	}
	task2 := &Task{
		Id:      2,
		running: 0,
		Name:    "task2",
	}
	taskList.Set("task1", task1)
	taskList.Set("task2", task2)
	taskList.EnqueuePending("task1")
	taskList.EnqueuePending("task2")
	t.Run("succuss", func(t *testing.T) {
		if tasks := taskList.GetAllPending(); len(tasks) != 2 {
			t.Errorf("GetAllPending() failed")
		}
	})
}

//杀死没有工作的任务
func Test_KillIdleTasks(t *testing.T) {
	taskList := &taskList{
		data: make(map[string]*TaskWithPending),
	}
	task1 := &Task{
		Id:      1,
		running: 0,
		Name:    "task1",
	}
	task2 := &Task{
		Id:      2,
		running: 0,
		Name:    "task2",
	}
	taskList.Set("task1", task1)
	taskList.Set("task2", task2)
	t.Run("succuss", func(t *testing.T) {
		taskList.KillIdleTasks()
		if len(taskList.data) != 0 {
			t.Errorf("KillIdleTasks() failed")
		}
	})
}

func Test_Len(t *testing.T) {
	taskList := &taskList{
		data: make(map[string]*TaskWithPending),
	}
	task1 := &Task{
		Id:      1,
		running: 0,
		Name:    "task1",
	}
	task2 := &Task{
		Id:      2,
		running: 0,
		Name:    "task2",
	}
	taskList.Set("task1", task1)
	taskList.Set("task2", task2)
	t.Run("succuss", func(t *testing.T) {
		if len := taskList.Len(); len != 2 {
			t.Errorf("Len() failed")
		}
	})
}

func Test_Del(t *testing.T) {
	taskList := &taskList{
		data: make(map[string]*TaskWithPending),
	}
	task1 := &Task{
		Id:      1,
		running: 0,
		Name:    "task1",
	}
	task2 := &Task{
		Id:      2,
		running: 0,
		Name:    "task2",
	}
	taskList.Set("task1", task1)
	taskList.Set("task2", task2)
	t.Run("succuss", func(t *testing.T) {
		taskList.Del("task1")
		got := taskList.Get("task1")
		if got != nil {
			t.Errorf("Del() failed")
		}
	})
}

func Test_Exists(t *testing.T) {
	taskList := &taskList{
		data: make(map[string]*TaskWithPending),
	}
	task1 := &Task{
		Id:      1,
		running: 0,
		Name:    "task1",
	}
	task2 := &Task{
		Id:      2,
		running: 0,
		Name:    "task2",
	}
	taskList.Set("task1", task1)
	taskList.Set("task2", task2)
	t.Run("succuss", func(t *testing.T) {
		if got := taskList.Exists("task1"); got != true {
			t.Errorf("Del() failed")
		}
	})
	t.Run("fail", func(t *testing.T) {
		if got := taskList.Exists("task3"); got != false {
			t.Errorf("Del() failed")
		}
	})
}
