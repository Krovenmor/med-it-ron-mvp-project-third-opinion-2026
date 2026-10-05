package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

const maxBodyBytes = 1 << 20

func (h *Handler) decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(dst); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "malformed JSON body"})
		return false
	}
	return true
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.log.Warn("write response", zap.Error(err))
	}
}
