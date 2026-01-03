package main

import (
	"fmt"
	_ "github.com/TubagusAldiMY/finance-tracker-app/backend/docs"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title           Finance Tracker API
// @version         1.0
// @description     API Documentation for Finance Tracker App.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Tubagus Aldi
// @contact.email   contact@tubsamy.tech

// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	// 3. Start Server
	webPort := viperConfig.GetInt("web.port")
	log.Infof("Server starting at port %d", webPort)
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
