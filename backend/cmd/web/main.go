package main

import (
	"fmt"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra"
)

func main() {
	// 1. Init Infrastructure
	viperConfig := infra.NewViper()
	log := infra.NewLogger(viperConfig)
	db := infra.NewDatabase(viperConfig, log)
	validate := infra.NewValidator(viperConfig)
	app := infra.NewFiber(viperConfig)

	// 2. Bootstrap Application (Wiring semua module di sini)
	infra.Bootstrap(&infra.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	// 3. Start Server
	webPort := viperConfig.GetInt("web.port")
	log.Infof("Server starting at port %d", webPort)
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
