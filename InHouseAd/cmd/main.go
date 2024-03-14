package main

import (
	"InHouseAd/internal/app"
	"InHouseAd/internal/config"
)

// @title User-Server
// @version 1.0
// @description User Server for InHouseAd test task.

// @host localhost:7777
// @BasePath /

func main() {
	cfg := config.NewConfig()
	app.Run(cfg)
}
