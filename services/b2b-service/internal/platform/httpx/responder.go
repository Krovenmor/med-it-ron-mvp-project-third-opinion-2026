package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

const (
	maxBodyBytes = 1 << 20
	ActorHeader  = "X-User-ID"
)

type Responder struct {
	log *zap.Logger
}

func NewResponder(log *zap.Logger, name string) Responder {
	return Responder{log: log.Named(name).Named("http")}
}

func (h Responder) PathID(w http.ResponseWriter, r *http.Request, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue(name))
	if err != nil {
		h.WriteError(w, r, apperr.ErrNotFound)
		return uuid.Nil, false
	}
	return id, true
}

func (h Responder) Actor(w http.ResponseWriter, r *http.Request) (string, bool) {
	actor := r.Header.Get(ActorHeader)
	if actor == "" {
		h.WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: ActorHeader + " header is required"})
		return "", false
	}
	return actor, true
}

func (h Responder) Decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(dst); err != nil {
		h.WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "malformed JSON body"})
		return false
	}
	return true
}

func (h Responder) DecodeOptional(w http.ResponseWriter, r *http.Request, dst any) bool {
	err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(dst)
	if err != nil && !errors.Is(err, io.EOF) {
		h.WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "malformed JSON body"})
		return false
	}
	return true
}

func (h Responder) WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.log.Warn("write response", zap.Error(err))
	}
}
