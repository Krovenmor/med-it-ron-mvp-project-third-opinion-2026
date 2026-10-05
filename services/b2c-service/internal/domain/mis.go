package domain

import "time"

type History struct {
	Visits       []Visit
	Appointments []Appointment
}

type Visit struct {
	ServiceCode string
	ServiceName string
	VisitedAt   time.Time
}

type Appointment struct {
	ID          string
	ServiceCode string
	ServiceName string
	ScheduledAt time.Time
	ReferralID  string
}

type Slot struct {
	ID          string
	ServiceCode string
	StartsAt    time.Time
}

type AppointmentRequest struct {
	PatientID   string
	SlotID      string
	ServiceName string
	ReferralID  string
}

type Booking struct {
	RecommendationID string
	AppointmentID    string
	ScheduledAt      time.Time
}

func (h History) AppointmentByReferral(recommendationID string) (Appointment, bool) {
	for _, a := range h.Appointments {
		if a.ReferralID == recommendationID {
			return a, true
		}
	}
	return Appointment{}, false
}
