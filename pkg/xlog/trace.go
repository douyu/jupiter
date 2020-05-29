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

package xlog

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

var _mdTrace = "_meta_trace"

// Tracer ...
type Tracer struct {
	BeginTime time.Time
	fields    []zap.Field
	mu        sync.RWMutex
	lv        zap.AtomicLevel
}

// NewTracer ...
func NewTracer() *Tracer {
	return &Tracer{
		BeginTime: time.Now(),
		fields:    make([]zap.Field, 0),
		lv:        zap.NewAtomicLevelAt(zap.InfoLevel),
	}
}

// Flush ...
func (t *Tracer) Flush(msg string, logger *Logger) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	switch t.lv.Level() {
	case zap.InfoLevel:
		logger.Info(msg, t.fields...)
	case zap.WarnLevel:
		logger.Warn(msg, t.fields...)
	case zap.ErrorLevel:
		logger.Error(msg, t.fields...)
	default:
	}
}

// Info ...
func (t *Tracer) Info(fields ...Field) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fields = append(t.fields, fields...)
	if t.lv.Level() < zap.InfoLevel {
		t.lv.SetLevel(zap.InfoLevel)
	}
}

// Warn ...
func (t *Tracer) Warn(fields ...Field) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fields = append(t.fields, fields...)
	if t.lv.Level() < zap.WarnLevel {
		t.lv.SetLevel(zap.WarnLevel)
	}
}

// Error ...
func (t *Tracer) Error(fields ...Field) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fields = append(t.fields, fields...)
	if t.lv.Level() < zap.ErrorLevel {
		t.lv.SetLevel(zap.ErrorLevel)
	}
}

// ExtractTraceMD ...
func ExtractTraceMD(ctx interface{ Get(string) interface{} }) (md *Tracer, ok bool) {
	md, ok = ctx.Get(_mdTrace).(*Tracer)
	return
}

// InjectTraceMD ...
func InjectTraceMD(ctx interface{ Set(string, interface{}) }, md *Tracer) {
	ctx.Set(_mdTrace, md)
}

type tracerKey struct{}

// NewContext ...
func NewContext(ctx context.Context, tracer Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}
