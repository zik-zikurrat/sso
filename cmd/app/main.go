package main

import (
	"sso/internal/app"
	"sso/internal/config"
)

func main() {
	// Config
	cfg := config.MustLoad()

	// Run application
	if err := app.Run(cfg); err != nil {
		return
	}
}
