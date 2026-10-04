package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

const (
	maxBodyBytes = 1 << 20
	actorHeader  = "X-User-ID"
)

func (h *Handler) pathID(w http.ResponseWriter, r *http.Request, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue(name))
	if err != nil {
		h.writeError(w, r, domain.ErrNotFound)
		return uuid.Nil, false
	}
	return id, true
}

func (h *Handler) actor(w http.ResponseWriter, r *http.Request) (string, bool) {
	actor := r.Header.Get(actorHeader)
	if actor == "" {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: actorHeader + " header is required"})
		return "", false
	}
	return actor, true
}

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
