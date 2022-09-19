package main

import (
	"api/config"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info(cfg.Log.Level)
}
