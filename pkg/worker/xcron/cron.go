// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xcron

import (
	"sync/atomic"
	"time"

	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/robfig/cron/v3"
)

var (
	// Every ...
	Every = cron.Every
	// NewParser ...
	NewParser = cron.NewParser
	// NewChain ...
	NewChain = cron.NewChain
	// WithSeconds ...
	WithSeconds = cron.WithSeconds
	// WithParser ...
	WithParser = cron.WithParser
	// WithLocation ...
	WithLocation = cron.WithLocation
)

type (
	// JobWrapper ...
	JobWrapper = cron.JobWrapper
	// EntryID ...
	EntryID = cron.EntryID
	// Entry ...
	Entry = cron.Entry
	// Schedule ...
	Schedule = cron.Schedule
	// Parser ...
	Parser = cron.Parser
	// Option ...
	Option = cron.Option
	// Job ...
	Job = cron.Job
	//NamedJob ..
	NamedJob interface {
		Run() error
		Name() string
	}
)

// FuncJob ...
type FuncJob func() error

// Run ...
func (f FuncJob) Run() error { return f() }

// Name ...
func (f FuncJob) Name() string { return xstring.FunctionName(f) }

// Cron ...
type Cron struct {
	*Config
	*cron.Cron
	entries map[string]EntryID
}

func newCron(config *Config) *Cron {
	if config.logger == nil {
		config.logger = xlog.JupiterLogger
	}
	config.logger = config.logger.With(xlog.FieldMod("worker.cron"))
	cron := &Cron{
		Config: config,
		Cron: cron.New(
			cron.WithParser(config.parser),
			cron.WithChain(config.wrappers...),
			cron.WithLogger(&wrappedLogger{config.logger}),
		),
	}
	return cron
}

// Schedule ...
func (c *Cron) Schedule(schedule Schedule, job NamedJob) EntryID {
	if c.ImmediatelyRun {
		schedule = &immediatelyScheduler{
			Schedule: schedule,
		}
	}
	innnerJob := &wrappedJob{
		NamedJob: job,
		logger:   c.logger,

		distributedTask: c.DistributedTask,
		waitLockTime:    c.WaitLockTime,
		leaseTTL:        c.Config.TTL,
		client:          c.client,
	}
	// xdebug.PrintKVWithPrefix("worker", "add job", job.Name())
	c.logger.Info("add job", xlog.String("name", job.Name()))
	return c.Cron.Schedule(schedule, innnerJob)
}

// GetEntryByName ...
func (c *Cron) GetEntryByName(name string) cron.Entry {
	// todo(gorexlv): data race
	return c.Entry(c.entries[name])
}

// AddJob ...
func (c *Cron) AddJob(spec string, cmd NamedJob) (EntryID, error) {
	schedule, err := c.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.Schedule(schedule, cmd), nil
}

// AddFunc ...
func (c *Cron) AddFunc(spec string, cmd func() error) (EntryID, error) {
	return c.AddJob(spec, FuncJob(cmd))
}

// Run ...
func (c *Cron) Run() error {
	// xdebug.PrintKVWithPrefix("worker", "run worker", fmt.Sprintf("%d job scheduled", len(c.Cron.Entries())))
	c.logger.Info("run worker", xlog.Int("number of scheduled jobs", len(c.Cron.Entries())))
	c.Cron.Run()
	return nil
}

// Stop ...
func (c *Cron) Stop() error {
	_ = c.Cron.Stop()
	return nil
}

type immediatelyScheduler struct {
	Schedule
	initOnce uint32
}

// Next ...
func (is *immediatelyScheduler) Next(curr time.Time) (next time.Time) {
	if atomic.CompareAndSwapUint32(&is.initOnce, 0, 1) {
		return curr
	}

	return is.Schedule.Next(curr)
}
