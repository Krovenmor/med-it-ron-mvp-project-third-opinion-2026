package http

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/scheduling"
)

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, scheduling.ErrInvalidRequest):
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, scheduling.ErrSlotTaken):
		h.writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	default:
		h.log.Error("request failed", zap.Error(err))
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
	}
}
