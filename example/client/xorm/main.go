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


