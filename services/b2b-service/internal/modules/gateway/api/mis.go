package api

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MIS interface {
	PatientHistory(ctx context.Context, patientID uuid.UUID) (History, error)
	Slots(ctx context.Context, serviceCode string, from time.Time) ([]Slot, error)
	Book(ctx context.Context, req AppointmentRequest) (Appointment, error)
}

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
	PatientID   uuid.UUID
	SlotID      string
	ServiceName string
	ReferralID  string
}
