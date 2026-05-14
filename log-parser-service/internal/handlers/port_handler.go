package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Jeno7u/log-parser/internal/dto"
	"github.com/Jeno7u/log-parser/internal/repository"
)

type PortHandler struct {
	repo repository.PortRepository
	log  *slog.Logger
}

func NewPort(repo repository.PortRepository, log *slog.Logger) *PortHandler {
	return &PortHandler{repo: repo, log: log}
}

func (h *PortHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "method not allowed"})
		return
	}

	nodeID := strings.TrimPrefix(r.URL.Path, "/api/v1/port/")
	if nodeID == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "node_id is required"})
		return
	}

	result, err := h.repo.ListPortsByNodeID(r.Context(), nodeID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}
