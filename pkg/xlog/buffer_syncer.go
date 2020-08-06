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
	"bufio"
	"context"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"
)

type bufferWriterSyncer struct {
	sync.Mutex
	bufferWriter *bufio.Writer
	ticker       *time.Ticker
}

const (
	// defaultBufferSize sizes the buffer associated with each WriterSync.
	defaultBufferSize = 256 * 1024

	// defaultFlushInterval means the default flush interval
	defaultFlushInterval = 30 * time.Second
)

// CloseFunc should be called when the caller exits to clean up buffers.
type CloseFunc func() error

// Buffer wraps a WriteSyncer in a buffer to improve performance,
// if bufferSize = 0, we set it to defaultBufferSize
// if flushInterval = 0, we set it to defaultFlushInterval
func Buffer(ws zapcore.WriteSyncer, bufferSize int, flushInterval time.Duration) (zapcore.WriteSyncer, CloseFunc) {
	if _, ok := ws.(*bufferWriterSyncer); ok {
		// no need to layer on another buffer
		return ws, func() error { return nil }
	}

	ctx, cancel := context.WithCancel(context.Background())

	if bufferSize == 0 {
		bufferSize = defaultBufferSize
	}

	if flushInterval == 0 {
		flushInterval = defaultFlushInterval
	}

	ticker := time.NewTicker(flushInterval)

	ws = &bufferWriterSyncer{
		bufferWriter: bufio.NewWriterSize(ws, bufferSize),
		ticker:       ticker,
	}

	// flush buffer every interval
	// we do not need to exit this goroutine until closefunc called explicitly
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// the background goroutine just keep syncing
				// until the close func is called.
				_ = ws.Sync()
			case <-ctx.Done():
				return
			}
		}
	}()

	closefunc := func() error {
		cancel()
		return ws.Sync()
	}

	return ws, closefunc
}

// Write ...
func (s *bufferWriterSyncer) Write(bs []byte) (int, error) {
	// bufio is not goroutine safe, so add lock writer here
	s.Lock()
	defer s.Unlock()

	// there are some logic internal for bufio.Writer here:
	// 1. when the buffer is enough, data would not be flushed.
	// 2. when the buffer is not enough, data would be flushed as soon as the buffer fills up.
	// this would lead to log spliting, which is not acceptable for log collector
	// so we need to flush bufferWriter before writing the data into bufferWriter
	if len(bs) > s.bufferWriter.Available() && s.bufferWriter.Buffered() > 0 {
		if err := s.bufferWriter.Flush(); err != nil {
			return 0, err
		}
	}

	return s.bufferWriter.Write(bs)
}

// Sync ...
func (s *bufferWriterSyncer) Sync() error {
	// bufio is not goroutine safe, so add lock writer here
	s.Lock()
	defer s.Unlock()

	return s.bufferWriter.Flush()
}
