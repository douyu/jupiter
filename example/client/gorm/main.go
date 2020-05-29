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
	"github.com/douyu/jupiter/pkg/store/gorm"
	"github.com/douyu/jupiter/pkg/xlog"
)

/**
1.新建一个数据库叫test
2.执行以下example，go run main.go --config=config.toml
*/
type User struct {
	Id   int    `gorm:"not null" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

func main() {
	eng := &jupiter.Application{}
	err := eng.Startup(
		func() error {
			gormDB := gorm.StdConfig("test").Build()
			models := []interface{}{
				&User{},
			}
			gormDB.SingularTable(true)
			gormDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models...)
			gormDB.Create(&User{
				Name: "jupiter",
			})

			var user User
			gormDB.Where("id = 1").Find(&user)
			xlog.Info("user info", xlog.String("name", user.Name))
			return nil
		},
	)
	if err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
}
