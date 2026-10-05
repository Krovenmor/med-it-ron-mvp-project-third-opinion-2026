package patient

import (
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
)

const actorPatient = "patient"

type Decline struct {
	CaseID           uuid.UUID
	RecommendationID uuid.UUID
	Reason           domain.DeclineReason
	Comment          string
}

type DeclineResult struct {
	Item       domain.PlanItem
	CaseStatus domain.RouteStatus
}

type HelpResult struct {
	Task    domain.OperatorTask
	Created bool
}
