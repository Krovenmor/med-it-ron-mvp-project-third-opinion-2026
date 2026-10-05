package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/operator"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

func (h *Handler) outcome(w http.ResponseWriter, r *http.Request) (operator.Outcome, bool) {
	actor, taskID, ok := h.taskAction(w, r)
	if !ok {
		return operator.Outcome{}, false
	}
	var req outcomeRequest
	if !h.DecodeOptional(w, r, &req) {
		return operator.Outcome{}, false
	}
	return req.toCommand(taskID, actor), true
}

func (h *Handler) taskAction(w http.ResponseWriter, r *http.Request) (string, uuid.UUID, bool) {
	actor, ok := h.Actor(w, r)
	if !ok {
		return "", uuid.Nil, false
	}
	taskID, ok := h.PathID(w, r, "id")
	return actor, taskID, ok
}

func (h *Handler) respondTask(w http.ResponseWriter, r *http.Request, task domain.OperatorTask, err error) {
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newOperatorTaskResponse(task))
}

func dashboardDays(raw string) (int, error) {
	days, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: days must be an integer", apperr.ErrInvalidInput)
	}
	return days, nil
}
