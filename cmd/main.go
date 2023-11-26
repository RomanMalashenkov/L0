package main

import (
	"fmt"
	"os"
	"test_wb/config"
	"test_wb/internal/app"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		fmt.Printf("Error at startup: %v", err)
		os.Exit(1)
	}

	app.Start(cfg)
}
