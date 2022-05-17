package main

import (
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/flag"
	"uuid/internal/app/uuidserver/server"
)

func main() {
	_ = flag.Parse()
	app := jupiter.DefaultApp()

	if err := server.InitApp(app); err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
