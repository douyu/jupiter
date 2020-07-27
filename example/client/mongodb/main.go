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
	"context"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/client/mongodb"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// run: go run main.go -config=config.toml
type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.exampleMongoDBInsertOne,
	); err != nil {
		xlog.Panic("startup mongodb", xlog.Any("err", err))
	}
	return eng
}
func main() {
	app := NewEngine()
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (eng *Engine) exampleMongoDBInsertOne() (err error) {
	//build client
	client := mongodb.StdConfig("mymongo").Build()
	var doc = bson.M{"_id": primitive.NewObjectID(), "hometown": "Atlanta"}
	ctx := context.Background()
	defer client.Disconnect(ctx)
	collection := client.Database("hawk").Collection("user")
	result, err := collection.InsertOne(ctx, doc)

	xlog.Info("insert one success", xlog.Any("insertID", result.InsertedID))
	return
}
