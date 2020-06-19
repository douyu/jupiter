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
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type errorWriter struct{}

func (*errorWriter) Write([]byte) (int, error) { return 0, errors.New("unimplemented") }

func requireWriteWorks(t testing.TB, ws zapcore.WriteSyncer) {
	n, err := ws.Write([]byte("foo"))
	require.NoError(t, err, "Unexpected error writing to WriteSyncer.")
	require.Equal(t, 3, n, "Wrote an unexpected number of bytes.")
}

func TestBufferWriter(t *testing.T) {
	// If we pass a plain io.Writer, make sure that we still get a WriteSyncer
	// with a no-op Sync.
	t.Run("sync", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ws, close := Buffer(zapcore.AddSync(buf), 0, 0)
		defer close()
		requireWriteWorks(t, ws)
		assert.Equal(t, "", buf.String(), "Unexpected log calling a no-op Write method.")
		assert.NoError(t, ws.Sync(), "Unexpected error calling a no-op Sync method.")
		assert.Equal(t, "foo", buf.String(), "Unexpected log string")
	})

	t.Run("1 close", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ws, close := Buffer(zapcore.AddSync(buf), 0, 0)
		requireWriteWorks(t, ws)
		assert.Equal(t, "", buf.String(), "Unexpected log calling a no-op Write method.")
		close()
		assert.Equal(t, "foo", buf.String(), "Unexpected log string")
	})

	t.Run("2 close", func(t *testing.T) {
		buf := &bytes.Buffer{}
		bufsync, close1 := Buffer(zapcore.AddSync(buf), 0, 0)
		ws, close2 := Buffer(bufsync, 0, 0)
		requireWriteWorks(t, ws)
		assert.Equal(t, "", buf.String(), "Unexpected log calling a no-op Write method.")
		close2()
		close1()
		assert.Equal(t, "foo", buf.String(), "Unexpected log string")
	})

	t.Run("small buffer", func(t *testing.T) {
		buf := &bytes.Buffer{}
		bufsync, close1 := Buffer(zapcore.AddSync(buf), 5, 0)
		ws, close2 := Buffer(bufsync, 5, 0)
		defer close1()
		defer close2()
		requireWriteWorks(t, ws)
		assert.Equal(t, "", buf.String(), "Unexpected log calling a no-op Write method.")
		requireWriteWorks(t, ws)
		assert.Equal(t, "foo", buf.String(), "Unexpected log string")
	})

	t.Run("flush error", func(t *testing.T) {
		ws, close := Buffer(zapcore.AddSync(&errorWriter{}), 4, time.Nanosecond)
		n, err := ws.Write([]byte("foo"))
		require.NoError(t, err, "Unexpected error writing to WriteSyncer.")
		require.Equal(t, 3, n, "Wrote an unexpected number of bytes.")
		ws.Write([]byte("foo"))
		assert.NotNil(t, close())
	})

	t.Run("flush timer", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ws, close := Buffer(zapcore.AddSync(buf), 6, time.Microsecond)
		defer close()
		requireWriteWorks(t, ws)
		time.Sleep(10 * time.Millisecond)
		bws := ws.(*bufferWriterSyncer)
		bws.Lock()
		assert.Equal(t, "foo", buf.String(), "Unexpected log string")
		bws.Unlock()

		// flush twice to validate loop logic
		requireWriteWorks(t, ws)
		time.Sleep(10 * time.Millisecond)
		bws.Lock()
		assert.Equal(t, "foofoo", buf.String(), "Unexpected log string")
		bws.Unlock()
	})
}
