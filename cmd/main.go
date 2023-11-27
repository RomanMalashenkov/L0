package main

import (
	"log"
	"test_wb/config"
	"test_wb/internal/app"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Error at startup: %v", err)
	}

	app.Start(cfg)
}
