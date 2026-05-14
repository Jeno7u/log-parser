package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Jeno7u/log-parser/internal/dto"
	"github.com/Jeno7u/log-parser/internal/repository"
)

type LogHandler struct {
	repo repository.LogRepository
	log  *slog.Logger
}

func NewLog(repo repository.LogRepository, log *slog.Logger) *LogHandler {
	return &LogHandler{repo: repo, log: log}
}

func (h *LogHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "method not allowed"})
		return
	}

	logID := strings.TrimPrefix(r.URL.Path, "/api/v1/log/")
	if logID == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "log_id is required"})
		return
	}

	result, err := h.repo.GetLogByID(r.Context(), logID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}
