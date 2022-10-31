package xxl

import (
	"errors"
	"sync"
)

var (
	OverLimit = "over limit"
)

type TaskWithPending struct {
	*Task
	pending chan *Task // 任务队列
}

//任务列表 [JobID]执行函数,并行执行时[+LogID]
type taskList struct {
	mu   sync.RWMutex
	data map[string]*TaskWithPending
}

//设置数据
func (t *taskList) Set(key string, val *Task) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.data[key] = &TaskWithPending{val, make(chan *Task, MaxQueueSize)}
}

//获取任务
func (t *taskList) Get(key string) *Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.data[key] != nil {
		return t.data[key].Task
	}
	return nil
}

//任务排入队列
func (t *taskList) EnqueuePending(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.data[key] != nil {
		l := len(t.data[key].pending)
		if l >= MaxQueueSize {
			return errors.New(OverLimit)
		}
		t.data[key].pending <- t.data[key].Task
		return nil
	}
	return errors.New("invalid key")
}

//获取队列
func (t *taskList) GetPending(key string) (*Task, chan *Task) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.data[key] != nil {
		return t.data[key].Task, t.data[key].pending
	}
	return nil, nil
}

//获取全部队列
func (t *taskList) GetAllPending() []*TaskWithPending {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var r []*TaskWithPending
	for _, task := range t.data {
		if len(task.pending) > 0 {
			r = append(r, task)
		}
	}
	return r
}

//杀死没有工作的任务
func (t *taskList) KillIdleTasks() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for key, task := range t.data {
		if len(task.pending) == 0 && !task.IsRunning() {
			close(task.pending)
			delete(t.data, key)
		}
	}
}

func (t *taskList) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.data)
}

//设置数据
func (t *taskList) Del(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.data[key] != nil {
		close(t.data[key].pending)
	}
	delete(t.data, key)
}

//Key是否存在
func (t *taskList) Exists(key string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.data[key]
	return ok
}
