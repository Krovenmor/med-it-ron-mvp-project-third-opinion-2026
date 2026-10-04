package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.log.Warn("write response", zap.Error(err))
	}
}
