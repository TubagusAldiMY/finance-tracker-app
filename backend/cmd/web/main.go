package main

import (
	"fmt"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/user"
)

func main() {
	// 1. Init Infrastructure
	viperConfig := infra.NewViper()
	log := infra.NewLogger(viperConfig)
	db := infra.NewDatabase(viperConfig, log)
	validate := infra.NewValidator(viperConfig)
	app := infra.NewFiber(viperConfig)

	// 2. Init Modules (Dependency Injection)
	// Module: User
	userRepo := user.NewRepository(db)
	userUseCase := user.NewUseCase(userRepo, log, validate)
	userHandler := user.NewHandler(userUseCase)

	// Register Routes dari Module User
	userHandler.RegisterRoutes(app)

	// 3. Start Server
	webPort := viperConfig.GetInt("web.port")
	log.Infof("Server starting at port %d", webPort)
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
