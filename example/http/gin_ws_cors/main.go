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
	"log"
	"net/http"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/server/xgin"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, "Hello Gin")
	})
	//Upgrade to websocket

	server.Upgrade(xgin.WebSocketOptions("/ws", handleWebSocketConn, handleCheckOrigin))
	return eng.Serve(server)
}

func handleWebSocketConn(ws xgin.WebSocketConn, err error) {
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

//handleCheckOrigin 允许websocket跨域
//error:request origin not allowed by Upgrader.CheckOrigin
func handleCheckOrigin(ws *xgin.WebSocket) {
	ws.Upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
