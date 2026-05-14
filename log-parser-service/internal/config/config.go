package config

import (
	"log/slog"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	PostgresConnString string
	Port               string
	LogLevel           string
	DataDir            string
}

func InitConfig(log *slog.Logger) Config {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	return Config{
		PostgresConnString: getEnv("DATABASE_URL", "postgres://postgres:password@db:5432/log-db?sslmode=disable", log),
		Port:               getEnv("PORT", "8080", log),
		LogLevel:           getEnv("LOG_LEVEL", "info", log),
		DataDir:            getEnv("DATA_DIR", "/app/data", log),
	}
}

func getEnv(key, fallback string, log *slog.Logger) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if log != nil {
		log.Info("cant find env by key: %v, using: %v", key, fallback)
	}

	return fallback
}
