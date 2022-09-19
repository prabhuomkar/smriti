package main

import (
	"api/config"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	log.Print(cfg.Log.Level)
}
