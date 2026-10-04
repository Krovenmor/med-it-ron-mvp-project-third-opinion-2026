package http

import "net/http"

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
