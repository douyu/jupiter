package main

import (
	"github.com/douyu/jupiter"
	"uuid/internal/uuidserver/server"
)

func main() {
	app := jupiter.DefaultApp()

	if err := server.InitApp(app); err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
