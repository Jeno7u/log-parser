package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Jeno7u/log-parser/internal/config"
	"github.com/Jeno7u/log-parser/internal/db"
	"github.com/Jeno7u/log-parser/internal/handlers"
	apphttp "github.com/Jeno7u/log-parser/internal/http"
	"github.com/Jeno7u/log-parser/internal/repository"
	"github.com/Jeno7u/log-parser/internal/service"
)

func main() {
	// creating logger and config setup
	log := config.NewLogger()
	config := config.InitConfig(log)

	// db setup
	db.MigrateDatabase(config.PostgresConnString, log)
	pool := db.NewPostgresPool(config.PostgresConnString, log)
	defer pool.Close()

	// creating repositories
	logRepository := repository.NewLogRepository(pool)
	nodeRepository := repository.NewNodeRepository(pool)
	portRepository := repository.NewPortRepository(pool)
	parseService := service.NewParse(config.DataDir, logRepository, nodeRepository, portRepository)

	// creating handlers
	parseHandler := handlers.NewParse(parseService, log)
	topologyHandler := handlers.NewTopology(logRepository, nodeRepository, portRepository, log)
	logHandler := handlers.NewLog(logRepository, log)
	nodeHandler := handlers.NewNode(nodeRepository, log)
	portHandler := handlers.NewPort(portRepository, log)

	// creating server instance and running it
	server := apphttp.Setup(config, log, parseHandler, topologyHandler, logHandler, nodeHandler, portHandler)
	if err := server.Run(context.Background()); err != nil {
		log.Error("server stopped with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
