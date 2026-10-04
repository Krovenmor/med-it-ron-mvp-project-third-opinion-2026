package booking

import (
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Book struct {
	CaseID           uuid.UUID
	RecommendationID uuid.UUID
	AppointmentID    string
	ScheduledAt      time.Time
	Channel          domain.BookingChannel
}

type Result struct {
	Booking    domain.Booking
	CaseStatus domain.CaseStatus
	Created    bool
}
