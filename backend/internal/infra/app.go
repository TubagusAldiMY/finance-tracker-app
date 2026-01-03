package infra

import (
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra/middleware"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/budget"
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

	// Auth Middleware
	authMiddleware := middleware.AuthMiddleware(config.Config)

	// User Module
	userRepo := user.NewRepository(config.DB)
	userUseCase := user.NewUseCase(userRepo, config.Log, config.Validate, config.Config)
	userHandler := user.NewHandler(userUseCase)

	userHandler.RegisterRoutes(config.App, authMiddleware)

	budgetRepo := budget.NewRepository(config.DB)
	budgetUseCase := budget.NewUseCase(budgetRepo, config.Log, config.Validate)
	budgetHandler := budget.NewHandler(budgetUseCase)

	budgetHandler.RegisterRoutes(config.App, authMiddleware)
}
