package http

import "net/http"

func (h *Handler) ingestReport(w http.ResponseWriter, r *http.Request) {
	var req reportRequest
	if !h.Decode(w, r, &req) {
		return
	}
	report, err := req.toDomain()
	if err != nil {
		h.WriteError(w, r, err)
		return
	}

	result, err := h.intake.Ingest(r.Context(), report)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}

	status := http.StatusOK
	if result.Created {
		status = http.StatusCreated
	}
	h.WriteJSON(w, status, newReportRefResponse(result.Intake))
}

func (h *Handler) reportStatus(w http.ResponseWriter, r *http.Request) {
	caseID, ok := h.PathID(w, r, "case_id")
	if !ok {
		return
	}
	intake, err := h.intake.Status(r.Context(), caseID)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newReportStatusResponse(intake))
}
