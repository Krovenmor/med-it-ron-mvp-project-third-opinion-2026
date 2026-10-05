package httpx

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

func (h Responder) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, apperr.ErrInvalidInput):
		h.WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apperr.ErrNotFound):
		h.WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: "not found"})
	case errors.Is(err, apperr.ErrInvalidState), errors.Is(err, apperr.ErrConflict):
		h.WriteJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apperr.ErrUpstream):
		h.log.Error("upstream failed", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
		h.WriteJSON(w, http.StatusBadGateway, ErrorResponse{Error: "upstream service unavailable"})
	default:
		h.log.Error("request failed", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
		h.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
	}
}
