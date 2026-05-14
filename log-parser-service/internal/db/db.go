package db

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

func NewPostgresPool(postgresConnString string, log *slog.Logger) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(postgresConnString)
	if err != nil {
		log.Error("got error when tried to parse db conn string", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// TODO: move this constants to env config
	// connection pool configuration
	poolConfig.MaxConns = 20
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Error("got error when tried to create conn pool", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("established connection pool to postgres successfully")
	return pool
}

// migrate db
func MigrateDatabase(postgresConnString string, log *slog.Logger) {
	db, err := sql.Open("pgx", postgresConnString)
	if err != nil {
		log.Error("got error when tried to establish connection during migration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// verify connection
	if err := db.Ping(); err != nil {
		log.Error("got error when tried to establish connection during migration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// perform migration
	err = goose.SetDialect("postgres")
	if err != nil {
		log.Error("got error when tried to set goose dialect", slog.String("error", err.Error()))
		os.Exit(1)
	}

	migrationsPath := firstExistingPath(
		"internal/db/migrations",
		"../db/migrations",
		"./migrations",
	)

	if err := goose.Up(db, migrationsPath); err != nil {
		log.Error("got error during migration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("migration was performed successfully")

}

func firstExistingPath(paths ...string) string {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return filepath.Clean(paths[0])
}
