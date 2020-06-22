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

package main

import (
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/server/xgin"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gin-gonic/gin"
)

// @title Demo
// @version 0.0.1
// @description  Demo接口描述文档
// @BasePath /
func main() {
	eng := NewEngine()
	if err := eng.Run(); err != nil {
		xlog.Panic(err.Error())
	}
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.serveHTTP,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

// HTTP地址
func (eng *Engine) serveHTTP() error {
	server := xgin.StdConfig("http").Build()
	server.GET("/hello", Hello)
	return eng.Serve(server)
}

// @Summary Hello
// @Description 详细说明Hello方法功能
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /hello [get]
func Hello(ctx *gin.Context) {
	ctx.JSON(200, NewOkResponse())
	return
}

// Response
type Response struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息描述
}

func NewOkResponse() *Response {
	return &Response{
		Code:    0,
		Message: "succeed",
	}
}
