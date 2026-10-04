package domain

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CaseStatus string

const (
	CaseStatusDraft       CaseStatus = "draft"
	CaseStatusInReview    CaseStatus = "in_review"
	CaseStatusConfirmed   CaseStatus = "confirmed"
	CaseStatusNotified    CaseStatus = "notified"
	CaseStatusBooked      CaseStatus = "booked"
	CaseStatusCompleted   CaseStatus = "completed"
	CaseStatusDeclined    CaseStatus = "declined"
	CaseStatusUnreachable CaseStatus = "unreachable"
)

type Urgency string

const (
	UrgencyNormal    Urgency = "normal"
	UrgencyPlanned   Urgency = "planned"
	UrgencyPriority  Urgency = "priority"
	UrgencyEmergency Urgency = "emergency"
)

func (u Urgency) Valid() bool {
	switch u {
	case UrgencyNormal, UrgencyPlanned, UrgencyPriority, UrgencyEmergency:
		return true
	}
	return false
}

type Modality string

const (
	ModalityCT Modality = "CT"
	ModalityDX Modality = "DX"
	ModalityMG Modality = "MG"
)

func (m Modality) Valid() bool {
	switch m {
	case ModalityCT, ModalityDX, ModalityMG:
		return true
	}
	return false
}

type Study struct {
	ID          string
	Modality    Modality
	BodySite    string
	PerformedAt time.Time
}

type Report struct {
	SourceSystem string
	Study        Study
	Conclusion   string
	Patient      Patient
}

func (r Report) Validate() error {
	switch {
	case r.SourceSystem == "":
		return invalid("source_system is required")
	case r.Study.ID == "":
		return invalid("study.id is required")
	case !r.Study.Modality.Valid():
		return invalid("study.modality must be one of CT, DX, MG")
	case r.Study.PerformedAt.IsZero():
		return invalid("study.performed_at is required")
	case r.Conclusion == "":
		return invalid("conclusion is required")
	case r.Patient.ExternalID == "":
		return invalid("patient.id is required")
	case r.Patient.FullName == "":
		return invalid("patient.full_name is required")
	case r.Patient.BirthDate.IsZero():
		return invalid("patient.birth_date is required")
	case !r.Patient.Sex.Valid():
		return invalid("patient.sex must be one of male, female")
	}
	return nil
}

func (r Report) Fingerprint() []byte {
	h := sha256.New()
	for _, part := range []string{
		r.Patient.ExternalID,
		string(r.Study.Modality),
		r.Study.BodySite,
		r.Study.PerformedAt.UTC().Format(time.RFC3339Nano),
		r.Conclusion,
	} {
		h.Write([]byte(part))
		h.Write([]byte{0})
	}
	return h.Sum(nil)
}

type Case struct {
	ID                uuid.UUID
	PatientID         uuid.UUID
	SourceSystem      string
	Study             Study
	Conclusion        string
	Fingerprint       []byte
	Status            CaseStatus
	Urgency           Urgency
	CatalogVersion    string
	GuidelinesVersion string
	ReceivedAt        time.Time
	UpdatedAt         time.Time
}

func NewCase(patientID uuid.UUID, r Report, now time.Time) Case {
	return Case{
		PatientID:    patientID,
		SourceSystem: r.SourceSystem,
		Study:        r.Study,
		Conclusion:   r.Conclusion,
		Fingerprint:  r.Fingerprint(),
		Status:       CaseStatusDraft,
		ReceivedAt:   now,
		UpdatedAt:    now,
	}
}

func (c *Case) ApplyAssessment(a Assessment, now time.Time) {
	c.Urgency = a.Urgency
	c.CatalogVersion = a.CatalogVersion
	c.GuidelinesVersion = a.GuidelinesVersion
	c.Status = CaseStatusInReview
	c.UpdatedAt = now
}

func invalid(reason string) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, reason)
}
