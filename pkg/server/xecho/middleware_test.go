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

package xecho

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware(t *testing.T) {
	e := echo.New()

	// Apply recovery middleware
	e.Use(recoveryMiddleware())

	// Handler that panics
	e.GET("/panic", func(c echo.Context) error {
		panic("test panic")
	})

	// Handler that returns error
	e.GET("/error", func(c echo.Context) error {
		panic(echo.NewHTTPError(http.StatusBadRequest, "bad request"))
	})

	// Normal handler
	e.GET("/normal", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	t.Run("should recover from panic", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Should not crash, error is recovered
		assert.NotEqual(t, http.StatusOK, rec.Code)
	})

	t.Run("should recover from error panic", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Should not crash
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("normal request should work", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/normal", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())
	})
}

func TestSlowLogMiddleware(t *testing.T) {
	e := echo.New()

	// Apply slow log middleware with 10ms threshold
	slowThreshold := 10 * time.Millisecond
	e.Use(slowLogMiddleware(slowThreshold))

	// Fast handler
	e.GET("/fast", func(c echo.Context) error {
		return c.String(http.StatusOK, "fast")
	})

	// Slow handler
	e.GET("/slow", func(c echo.Context) error {
		time.Sleep(20 * time.Millisecond)
		return c.String(http.StatusOK, "slow")
	})

	t.Run("fast request should complete normally", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/fast", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "fast", rec.Body.String())
	})

	t.Run("slow request should complete normally and log slow", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slow", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "slow", rec.Body.String())
	})
}

func TestSlowLogMiddleware_ZeroThreshold(t *testing.T) {
	e := echo.New()

	// Apply slow log middleware with 0 threshold (disabled)
	e.Use(slowLogMiddleware(0))

	e.GET("/test", func(c echo.Context) error {
		time.Sleep(5 * time.Millisecond)
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
