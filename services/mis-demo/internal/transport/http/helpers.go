package http

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.log.Warn("write response", zap.Error(err))
	}
}

func decodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func slotsFrom(raw string, now time.Time) (time.Time, bool) {
	if raw == "" {
		return now, true
	}
	from, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, false
	}
	if from.Before(now) {
		return now, true
	}
	return from, true
}
