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
	"github.com/douyu/jupiter/pkg/xlog"
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
	server := xwebsocket.StdConfig("websocket").Build()
	server.WebSocket("/ws", func(ctx *websocket.Conn, err error) {
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		for {
			mt, message, err := ctx.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			err = ctx.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
		return
	})
	return eng.Serve(server)
}
func echo(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{} // use default options
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

}
