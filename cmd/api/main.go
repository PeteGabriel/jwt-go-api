package main

import (
	"jwtGoApi/internal/api"
	"jwtGoApi/pkg/config"
	"jwtGoApi/pkg/data"
)

func main() {

	cfg := config.New()
	conn := data.NewConnection(cfg)
	defer conn.Disconnect()

	app := api.New(cfg, conn.Client)
	app.Start()
}
