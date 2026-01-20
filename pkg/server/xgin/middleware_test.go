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

package xgin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRecoveryMiddleware(t *testing.T) {
	logger := xlog.Jupiter()

	t.Run("should recover from string panic", func(t *testing.T) {
		r := gin.New()
		r.Use(recoveryMiddleware(logger))

		r.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("should recover from error panic", func(t *testing.T) {
		r := gin.New()
		r.Use(recoveryMiddleware(logger))

		r.GET("/error", func(c *gin.Context) {
			panic(gin.Error{Err: http.ErrAbortHandler})
		})

		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("normal request should work", func(t *testing.T) {
		r := gin.New()
		r.Use(recoveryMiddleware(logger))

		r.GET("/normal", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/normal", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())
	})
}

func TestSlowLogMiddleware(t *testing.T) {
	logger := xlog.Jupiter()

	t.Run("fast request should complete normally", func(t *testing.T) {
		r := gin.New()
		r.Use(slowLogMiddleware(logger, 10*time.Millisecond))

		r.GET("/fast", func(c *gin.Context) {
			c.String(http.StatusOK, "fast")
		})

		req := httptest.NewRequest(http.MethodGet, "/fast", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "fast", rec.Body.String())
	})

	t.Run("slow request should complete normally and log slow", func(t *testing.T) {
		r := gin.New()
		r.Use(slowLogMiddleware(logger, 10*time.Millisecond))

		r.GET("/slow", func(c *gin.Context) {
			time.Sleep(20 * time.Millisecond)
			c.String(http.StatusOK, "slow")
		})

		req := httptest.NewRequest(http.MethodGet, "/slow", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "slow", rec.Body.String())
	})
}

func TestSlowLogMiddleware_ZeroThreshold(t *testing.T) {
	logger := xlog.Jupiter()

	r := gin.New()
	r.Use(slowLogMiddleware(logger, 0))

	r.GET("/test", func(c *gin.Context) {
		time.Sleep(5 * time.Millisecond)
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
