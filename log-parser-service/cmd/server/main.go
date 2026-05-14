package main

import (
	"github.com/Jeno7u/log-parser/internal/config"
	"github.com/Jeno7u/log-parser/internal/db"
	"github.com/Jeno7u/log-parser/internal/repository"
)

func main() {
	// create config struct
	config := config.InitConfig()

	// setup db
	pool := db.NewPostgresPool(config.PostgresConnString)
	defer pool.Close()

	// create repository, service and other stuff
	logRepository := repository.NewLogRepository(pool)
}
