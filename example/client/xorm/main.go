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
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/store/xorm"
	"github.com/douyu/jupiter/pkg/xlog"
)

func main (){


	eng := &jupiter.Application{}
	if err := eng.Startup(openDB); err != nil {
		xlog.Panic("start up", xlog.FieldErr(err))
	}

	if err := eng.Run(); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
}

var Db xorm.Db

type User struct {
	Id int64 `xorm:"id"`
	Password string `xorm:"password"`
}

func openDB() error {

	Db = xorm.StdConfig("test").Build()

	_ = Db.Ping()

	var getUserOne User

	_, err := Db.Where("id = ?",101).Get(&getUserOne)

	fmt.Println(getUserOne,err)
	return nil
}


