package operator

import (
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
)

type Outcome struct {
	TaskID  uuid.UUID
	Actor   string
	Comment string
}

type Callback struct {
	Outcome
	At time.Time
}

type Decline struct {
	Outcome
	Reason domain.DeclineReason
}

type Book struct {
	Outcome
	RecommendationID uuid.UUID
	SlotID           string
}

type Card struct {
	Task     domain.OperatorTask
	Route    domain.Route
	Patient  domain.Patient
	Attempts []domain.CallAttempt
	Offers   []Offer
}

type Offer struct {
	Item             domain.PlanItem
	Booked           bool
	Slots            []gatewayapi.Slot
	SlotsUnavailable bool
}
