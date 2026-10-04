package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/cases"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/intake"
)

const maxBodyBytes = 1 << 20

type Intake interface {
	Ingest(ctx context.Context, report domain.Report) (intake.IngestResult, error)
}

type Cases interface {
	Get(ctx context.Context, id uuid.UUID) (cases.Details, error)
}

type DemoClock interface {
	Advance(d time.Duration) time.Time
}

type Handler struct {
	intake Intake
	cases  Cases
	clock  DemoClock
	log    *zap.Logger
}

func NewHandler(intake Intake, cases Cases, clock DemoClock, log *zap.Logger) *Handler {
	return &Handler{intake: intake, cases: cases, clock: clock, log: log.Named("http")}
}

func (h *Handler) Routes(demoMode bool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/reports", h.ingestReport)
	mux.HandleFunc("GET /api/v1/cases/{id}", h.getCase)
	if demoMode {
		mux.HandleFunc("POST /demo/clock/advance", h.advanceClock)
	}
	return h.recoverPanics(h.logRequests(mux))
}

func (h *Handler) ingestReport(w http.ResponseWriter, r *http.Request) {
	var req reportRequest
	if !h.decode(w, r, &req) {
		return
	}
	report, err := req.toDomain()
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	result, err := h.intake.Ingest(r.Context(), report)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	status := http.StatusOK
	if result.Created {
		status = http.StatusCreated
	}
	h.writeJSON(w, status, caseRefResponse{CaseID: result.Case.ID, Status: string(result.Case.Status)})
}

func (h *Handler) getCase(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.writeError(w, r, domain.ErrNotFound)
		return
	}
	details, err := h.cases.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusOK, newCaseResponse(details))
}

func (h *Handler) advanceClock(w http.ResponseWriter, r *http.Request) {
	var req advanceClockRequest
	if !h.decode(w, r, &req) {
		return
	}
	d, err := time.ParseDuration(req.By)
	if err != nil || d <= 0 {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: `"by" must be a positive duration, e.g. "30m" or "24h"`})
		return
	}
	h.writeJSON(w, http.StatusOK, clockResponse{Now: h.clock.Advance(d)})
}

func (h *Handler) decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(dst); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "malformed JSON body"})
		return false
	}
	return true
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "not found"})
	case errors.Is(err, domain.ErrStudyConflict):
		h.writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	default:
		h.log.Error("request failed", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
	}
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.log.Warn("write response", zap.Error(err))
	}
}
