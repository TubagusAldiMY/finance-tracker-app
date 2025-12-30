package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewDatabase(v *viper.Viper, log *logrus.Logger) *pgxpool.Pool {
	username := v.GetString("database.username")
	password := v.GetString("database.password")
	host := v.GetString("database.host")
	port := v.GetInt("database.port")
	database := v.GetString("database.name")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, database)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
	}

	config.MaxConns = int32(v.GetInt("database.pool.max"))
	config.MinConns = int32(v.GetInt("database.pool.idle"))
	config.MaxConnLifetime = time.Duration(v.GetInt("database.pool.lifetime")) * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Test Connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Info("Connected to DB successfully")
	return pool
}
