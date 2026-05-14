package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Jeno7u/log-parser/internal/config"
	"github.com/Jeno7u/log-parser/internal/handlers"
	"github.com/Jeno7u/log-parser/internal/middleware"
)

type Server struct {
	server *http.Server
	log    *slog.Logger
	config config.Config
}

func Setup(cfg config.Config, log *slog.Logger, parseHandler *handlers.ParseHandler, topologyHandler *handlers.TopologyHandler, logHandler *handlers.LogHandler, nodeHandler *handlers.NodeHandler, portHandler *handlers.PortHandler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/parse/", parseHandler.Handle)
	mux.HandleFunc("/api/v1/topology/", topologyHandler.Handle)
	mux.HandleFunc("/api/v1/log/", logHandler.Handle)
	mux.HandleFunc("/api/v1/node/", nodeHandler.Handle)
	mux.HandleFunc("/api/v1/port/", portHandler.Handle)

	handler := middleware.Recovery(log, middleware.RequestLogger(log, mux))

	return &Server{
		server: &http.Server{Addr: ":" + cfg.Port, Handler: handler, ReadHeaderTimeout: 5 * time.Second},
		log:    log,
		config: cfg,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = s.server.Shutdown(shutdownCtx)
	}()

	s.log.Info("http server started", slog.String("port", s.config.Port))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
