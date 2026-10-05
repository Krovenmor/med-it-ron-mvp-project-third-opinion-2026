package http

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "not found"})
	case errors.Is(err, domain.ErrInvalidState), errors.Is(err, domain.ErrSlotUnavailable):
		h.writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUpstream):
		h.log.Error("upstream failed", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
		h.writeJSON(w, http.StatusBadGateway, errorResponse{Error: "service temporarily unavailable"})
	default:
		h.log.Error("request failed", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
	}
}
