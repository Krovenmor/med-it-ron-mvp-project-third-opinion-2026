package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

type BookingChannel string

const (
	BookingChannelSelf     BookingChannel = "self"
	BookingChannelOperator BookingChannel = "operator"
)

func (c BookingChannel) Valid() bool {
	return c == BookingChannelSelf || c == BookingChannelOperator
}

func (c BookingChannel) Actor() string {
	if c == BookingChannelSelf {
		return "patient"
	}
	return string(c)
}

type Booking struct {
	ID               uuid.UUID
	CaseID           uuid.UUID
	RecommendationID uuid.UUID
	AppointmentID    string
	ScheduledAt      time.Time
	Channel          BookingChannel
	CreatedAt        time.Time
}

func (b Booking) Validate() error {
	switch {
	case b.AppointmentID == "":
		return apperr.Invalid("appointment_id is required")
	case b.ScheduledAt.IsZero():
		return apperr.Invalid("scheduled_at is required")
	case !b.Channel.Valid():
		return apperr.Invalid("channel must be one of self, operator")
	}
	return nil
}

func (b Booking) SameTarget(other Booking) bool {
	return b.CaseID == other.CaseID && b.RecommendationID == other.RecommendationID
}
