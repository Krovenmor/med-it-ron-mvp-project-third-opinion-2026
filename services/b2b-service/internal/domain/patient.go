package domain

import (
	"time"

	"github.com/google/uuid"
)

type Sex string

const (
	SexMale   Sex = "male"
	SexFemale Sex = "female"
)

func (s Sex) Valid() bool {
	return s == SexMale || s == SexFemale
}

type Patient struct {
	ID           uuid.UUID
	SourceSystem string
	ExternalID   string
	FullName     string
	BirthDate    time.Time
	Sex          Sex
	Phone        string
	Email        string
}

func (p Patient) AgeAt(t time.Time) int {
	age := t.Year() - p.BirthDate.Year()
	if t.Month() < p.BirthDate.Month() || (t.Month() == p.BirthDate.Month() && t.Day() < p.BirthDate.Day()) {
		age--
	}
	return age
}

type Visit struct {
	ServiceCode string
	ServiceName string
	VisitedAt   time.Time
}

type Appointment struct {
	ServiceCode string
	ServiceName string
	ScheduledAt time.Time
}

type PatientHistory struct {
	Visits       []Visit
	Appointments []Appointment
}
