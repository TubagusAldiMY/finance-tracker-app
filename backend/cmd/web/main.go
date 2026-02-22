package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/TubagusAldiMY/finance-tracker-app/backend/docs"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title           Finance Tracker API
// @version         1.0
// @description     API Documentation for Finance Tracker App.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Tubagus Aldi
// @contact.email   admin@tubsamy.tech

// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Init Infrastructure
	viperConfig := infra.NewViper("config")
	log := infra.NewLogger(viperConfig)
	db := infra.NewDatabase(viperConfig, log)
	validate := infra.NewValidator(viperConfig)
	app := infra.NewFiber(viperConfig)

	// 2. Bootstrap Application
	infra.Bootstrap(&infra.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	// 3. Security & Compliance: Swagger Toggle
	// Jangan ekspos dokumentasi internal di Production
	env := viperConfig.GetString("app.env") // Pastikan key ini ada di config.json Anda
	if env == "development" || env == "staging" {
		app.Get("/swagger/*", fiberSwagger.WrapHandler)
		log.Info("Swagger UI enabled at /swagger/index.html")
	}

	// 4. Start Server with Graceful Shutdown pattern
	webPort := viperConfig.GetInt("web.port")
	serverAddr := fmt.Sprintf(":%d", webPort)

	// Channel untuk mendengarkan sinyal OS
	// SIGINT (Ctrl+C), SIGTERM (Docker stop / K8s terminate)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Jalankan server di Goroutine terpisah agar tidak memblokir main thread
	go func() {
		log.Infof("Server starting at port %d in %s mode", webPort, env)
		if err := app.Listen(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 5. Blocking Main Thread sampai sinyal diterima
	<-quit
	log.Warn("Shutting down server...")

	// 6. Graceful Shutdown Timeout
	// Beri waktu (misal 10 detik) untuk menyelesaikan request yang sedang berjalan
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Matikan Fiber
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	// 7. Cleanup Resource Lain (Opsional tapi disarankan)
	// Tutup koneksi DB pool secara eksplisit
	log.Info("Closing database connection pool...")
	db.Close()

	log.Info("Server exited properly")
}
