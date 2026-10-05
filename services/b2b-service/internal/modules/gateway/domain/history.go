package domain

import "time"

type PatientHistory struct {
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
