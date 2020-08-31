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
	"fmt"
	"time"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/store/mongo"
	"github.com/douyu/jupiter/pkg/xlog"
	"gopkg.in/mgo.v2/bson"
)

// run: go run main.go -config=config.toml
type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.exampleMongo,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

func main() {
	app := NewEngine()
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (eng *Engine) exampleMongo() (err error) {
	session := mongo.StdNew("test")
	defer session.Close()

	// write
	m := make(map[string]interface{}, 0)
	m["dateline"] = time.Now().Unix()
	m["rid"] = 777

	err = session.DB("test").C("test").Insert(m)
	if err != nil {
		panic(err)
	}

	// read
	type MongoData struct {
		Dateline int64 `bson:"dateline"`
		Rid      int64 `bson:"rid"`
	}
	var rawData []MongoData
	err = session.DB("test").C("test").Find(bson.M{"rid": 777}).Sort("-time").All(&rawData)
	if err != nil {
		panic(err)
	}
	fmt.Println("rawData...", rawData)
	return
}
