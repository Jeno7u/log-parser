package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/Jeno7u/log-parser/internal/dto"
	"github.com/Jeno7u/log-parser/internal/service"
)

type ParseHandler struct {
	service service.ParseService
	log     *slog.Logger
}

func NewParse(service service.ParseService, log *slog.Logger) *ParseHandler {
	return &ParseHandler{service: service, log: log}
}

func (h *ParseHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "method not allowed"})
		return
	}

	var request dto.ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid json body"})
		return
	}
	if request.Path == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "path is required"})
		return
	}

	logID, err := h.service.Parse(context.Background(), request.Path, filepath.Base(request.Path))
	if err != nil {
		h.log.Error("parse failed", slog.String("error", err.Error()))
		writeJSON(w, http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusAccepted, dto.ParseResponse{LogID: logID})
}
