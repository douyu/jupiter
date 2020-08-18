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
	_ "github.com/douyu/jupiter/example/http/gin-swagger/docs"
	"github.com/douyu/jupiter/pkg/server/xgin"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1

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

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	server.GET("/v1/user/:id", getUser)
	return eng.Serve(server)
}

type Response struct {
	Code int
	Data interface{}
}

type User struct {
	Id   int
	Name string
}

// @Summary get user by ID
// @Description get user by ID
// @Accept  json
// @Produce json
// @Param id  path int true "user ID"
// @Success 200 object Response "success"
// @Success 400 object Response "failure"
// @Router /user/{id} [get]
func getUser(ctx *gin.Context) {
	user := User{
		Id:   1,
		Name: "xxx",
	}
	resp := Response{
		Code: 0,
		Data: user,
	}
	ctx.JSON(200, resp)
}
