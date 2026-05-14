package config

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger() *slog.Logger {
	// getting logger lever from env
	level := getEnv("LOG_LEVEL", "info", nil)

	// creating new instance of logger with needed level
	var slogLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn", "warning":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	return slog.New(handler)
}
