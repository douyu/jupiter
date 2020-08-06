package main

import (
	"fmt"
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/xlog"
)

// go run main.go --job=jobrunner
func main() {
	eng := NewEngine()
	if err := eng.Run(); err != nil {
		xlog.Error(err.Error())
	}
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.initJob,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

func (e *Engine) initJob() error {
	return e.Job(NewJobRunner())
}

type JobRunner struct {
	JobName string
}

func NewJobRunner() *JobRunner {
	return &JobRunner{
		JobName: "jobrunner",
	}
}

func (j *JobRunner) Run() {
	fmt.Println("i am job runner")
}

func (j *JobRunner) GetJobName() string {
	return j.JobName
}
