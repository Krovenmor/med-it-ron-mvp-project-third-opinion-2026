package domain

import (
	"time"

	"github.com/google/uuid"
)

type IntakeStatus string

const (
	IntakeStatusReceived IntakeStatus = "received"
	IntakeStatusAssessed IntakeStatus = "assessed"
)

type Intake struct {
	CaseID       uuid.UUID
	PatientID    uuid.UUID
	SourceSystem string
	Study        Study
	Conclusion   string
	Fingerprint  []byte
	Status       IntakeStatus
	ReceivedAt   time.Time
	AssessedAt   time.Time
}

func NewIntake(patientID uuid.UUID, r Report, now time.Time) Intake {
	return Intake{
		PatientID:    patientID,
		SourceSystem: r.SourceSystem,
		Study:        r.Study,
		Conclusion:   r.Conclusion,
		Fingerprint:  r.Fingerprint(),
		Status:       IntakeStatusReceived,
		ReceivedAt:   now,
	}
}

func (i *Intake) MarkAssessed(now time.Time) bool {
	if i.Status != IntakeStatusReceived {
		return false
	}
	i.Status = IntakeStatusAssessed
	i.AssessedAt = now
	return true
}
