package config

import (
	"fmt"
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	PostgresConnString string
	Port               string
	LogLevel           string
	DataDir            string
}

func InitConfig() Config {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	postgresConnString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		getEnv("POSTGRES_USER", "postgres"),
		getEnv("POSTGRES_PASSWORD", "password"),
		getEnv("POSTGRES_HOST", "postgres"),
		getEnv("POSTGRES_DB", "db"),
	)

	return Config{
		PostgresConnString: postgresConnString,
		Port:               getEnv("PORT", "8080"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		DataDir:            getEnv("DATA_DIR", "/app/data"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	log.Printf("cant find env by key: %v, using: %v", key, fallback)
	return fallback
}
