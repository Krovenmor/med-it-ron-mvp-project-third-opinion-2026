package domain

import "time"

type Urgency string

const (
	UrgencyNormal    Urgency = "normal"
	UrgencyPlanned   Urgency = "planned"
	UrgencyPriority  Urgency = "priority"
	UrgencyEmergency Urgency = "emergency"
)

const ReviewEscalationDelay = 30 * time.Minute

func (u Urgency) Valid() bool {
	return u.rank() >= 0
}

func (u Urgency) rank() int {
	switch u {
	case UrgencyNormal:
		return 0
	case UrgencyPlanned:
		return 1
	case UrgencyPriority:
		return 2
	case UrgencyEmergency:
		return 3
	}
	return -1
}

func (u Urgency) ReviewSLA() time.Duration {
	switch u {
	case UrgencyEmergency:
		return ReviewEscalationDelay
	case UrgencyPriority:
		return 4 * time.Hour
	default:
		return 24 * time.Hour
	}
}
