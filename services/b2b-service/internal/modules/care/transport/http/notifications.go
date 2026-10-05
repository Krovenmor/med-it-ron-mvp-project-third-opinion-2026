package http

import "net/http"

func (h *Handler) notificationJournal(w http.ResponseWriter, r *http.Request) {
	entries, err := h.operator.Journal(r.Context())
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newNotificationsResponse(entries))
}
