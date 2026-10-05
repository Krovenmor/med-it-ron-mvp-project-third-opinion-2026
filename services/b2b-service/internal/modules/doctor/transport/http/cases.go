package http

import "net/http"

func (h *Handler) getCase(w http.ResponseWriter, r *http.Request) {
	id, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	details, err := h.cases.Get(r.Context(), id)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newCaseResponse(details))
}

func (h *Handler) caseHistory(w http.ResponseWriter, r *http.Request) {
	id, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	history, err := h.cases.History(r.Context(), id)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newCaseHistoryResponse(history))
}
