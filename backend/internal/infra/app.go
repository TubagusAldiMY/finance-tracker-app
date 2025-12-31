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

	userRepo := user.NewRepository(config.DB)

	userUseCase := user.NewUseCase(userRepo, config.Log, config.Validate, config.Config)

	userHandler := user.NewHandler(userUseCase)

	userHandler.RegisterRoutes(config.App)

}
