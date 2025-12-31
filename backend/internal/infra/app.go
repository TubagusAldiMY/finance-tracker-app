package infra

import (
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/user" // Import module User

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB       *pgxpool.Pool
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	// 1. Setup Repository
	// Repository sekarang butuh *pgxpool.Pool langsung
	userRepo := user.NewRepository(config.DB)

	// 2. Setup UseCase
	// UseCase butuh Repository, Logger, dan Validator
	userUseCase := user.NewUseCase(userRepo, config.Log, config.Validate)

	// 3. Setup Handler (Controller)
	userHandler := user.NewHandler(userUseCase)

	// 4. Register Routes
	// Route didaftarkan langsung oleh handlernya (Self-contained)
	userHandler.RegisterRoutes(config.App)

}
