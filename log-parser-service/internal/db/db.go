package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

func NewPostgresPool(postgresConnString string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(postgresConnString)
	if err != nil {
		log.Fatalln("got error when tried to parse db conn string, ", err)
	}

	// TODO: move this constants to env config
	// connection pool configuration
	poolConfig.MaxConns = 20
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("got error when tried to create conn pool, ", err)
	}

	log.Println("Established connection pool to postgres. Success!")
	return pool
}

// migrate db
func MigrateDatabase(postgresConnString string) {
	db, err := sql.Open("pgx", postgresConnString)
	if err != nil {
		log.Fatalln("got error when tried to establish connection during migration, ", err)
	}

	// verify connection
	if err := db.Ping(); err != nil {
		log.Fatalln("got error when tried to establish connection during migration, ", err)
	}

	// perform migration
	err = goose.SetDialect("postgres")
	if err != nil {
		log.Fatalln("got error when tried to set goose dialect, ", err)
	}

	migrationsPath := firstExistingPath(
		"internal/db/migrations",
		"../db/migrations",
		"./migrations",
	)

	if err := goose.Up(db, migrationsPath); err != nil {
		log.Fatalln("got error during migration, ", err)
	}

	log.Println("migration was performed successfully!")

}

func firstExistingPath(paths ...string) string {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return filepath.Clean(paths[0])
}
