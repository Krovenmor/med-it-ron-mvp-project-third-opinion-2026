package domain

import (
	"crypto/sha256"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

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
		return apperr.Invalid("source_system is required")
	case r.Study.ID == "":
		return apperr.Invalid("study.id is required")
	case !r.Study.Modality.Valid():
		return apperr.Invalid("study.modality must be one of CT, DX, MG")
	case r.Study.PerformedAt.IsZero():
		return apperr.Invalid("study.performed_at is required")
	case r.Conclusion == "":
		return apperr.Invalid("conclusion is required")
	case r.Patient.ExternalID == "":
		return apperr.Invalid("patient.id is required")
	case r.Patient.FullName == "":
		return apperr.Invalid("patient.full_name is required")
	case r.Patient.BirthDate.IsZero():
		return apperr.Invalid("patient.birth_date is required")
	case !r.Patient.Sex.Valid():
		return apperr.Invalid("patient.sex must be one of male, female")
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
