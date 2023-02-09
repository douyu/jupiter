// Copyright 2022 Douyu
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
	"github.com/labstack/echo/v4"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Server(t *testing.T) {
	c := DefaultConfig()
	c.Port = 0
	s := c.MustBuild()
	s.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": "hello jupiter",
		})
	})
	s.GET("/api/jupiter/biz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": "hello jupiter",
		})
	})
	s.POST("/api/jupiter/internal", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": "hello jupiter",
		})
	})
	s.POST("/api/jupiter/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": "hello jupiter",
		})
	})
	go func() {
		s.Serve()
	}()
	time.Sleep(time.Second)
	assert.True(t, s.Healthz())
	assert.NotNil(t, s.Info())
	s.Stop()
}
