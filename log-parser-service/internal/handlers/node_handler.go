package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Jeno7u/log-parser/internal/dto"
	"github.com/Jeno7u/log-parser/internal/repository"
)

type NodeHandler struct {
	repo repository.NodeRepository
	log  *slog.Logger
}

func NewNode(repo repository.NodeRepository, log *slog.Logger) *NodeHandler {
	return &NodeHandler{repo: repo, log: log}
}

func (h *NodeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "method not allowed"})
		return
	}

	nodeID := strings.TrimPrefix(r.URL.Path, "/api/v1/node/")
	if nodeID == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "node_id is required"})
		return
	}

	result, err := h.repo.GetNodeByID(r.Context(), nodeID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}
