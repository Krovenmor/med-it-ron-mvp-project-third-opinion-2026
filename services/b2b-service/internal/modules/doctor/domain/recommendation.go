package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

type RecommendationSource string

const (
	RecommendationSourceAI     RecommendationSource = "ai"
	RecommendationSourceDoctor RecommendationSource = "doctor"
)

type Importance string

const (
	ImportanceHigh Importance = "high"
	ImportanceLow  Importance = "low"
)

func (i Importance) Valid() bool {
	return i == ImportanceHigh || i == ImportanceLow
}

type Mark string

const (
	MarkCritical Mark = "critical"
	MarkMinor    Mark = "minor"
	MarkRejected Mark = "rejected"
)

type RejectReason string

const (
	RejectReasonContraindicated RejectReason = "contraindicated"
	RejectReasonOther           RejectReason = "other"
)

type Review struct {
	Mark          Mark
	RejectReason  RejectReason
	RejectComment string
	ReviewedBy    string
	ReviewedAt    time.Time
}

func (r Review) Validate() error {
	switch r.Mark {
	case MarkCritical, MarkMinor:
		if r.RejectReason != "" || r.RejectComment != "" {
			return apperr.Invalid("reject_reason and reject_comment are allowed only for rejected recommendations")
		}
	case MarkRejected:
		switch r.RejectReason {
		case RejectReasonContraindicated:
		case RejectReasonOther:
			if r.RejectComment == "" {
				return apperr.Invalid("reject_comment is required when reject_reason is other")
			}
		default:
			return apperr.Invalid("reject_reason must be one of contraindicated, other")
		}
	default:
		return apperr.Invalid("mark must be one of critical, minor, rejected")
	}
	return nil
}

type Recommendation struct {
	ID            uuid.UUID
	CaseID        uuid.UUID
	Position      int
	Source        RecommendationSource
	ServiceCode   string
	ServiceName   string
	Importance    Importance
	Rationale     string
	GuidelineRef  string
	PatientText   string
	AlreadyBooked bool
	Review        Review
	CreatedAt     time.Time
}

func NewDoctorRecommendation(caseID uuid.UUID, service Service, patientText, rationale string, review Review) (Recommendation, error) {
	switch {
	case service.Name == "":
		return Recommendation{}, apperr.Invalid("service_name is required")
	case patientText == "":
		return Recommendation{}, apperr.Invalid("patient_text is required")
	case review.Mark != MarkCritical && review.Mark != MarkMinor:
		return Recommendation{}, apperr.Invalid("mark must be one of critical, minor")
	}

	importance := ImportanceLow
	if review.Mark == MarkCritical {
		importance = ImportanceHigh
	}
	return Recommendation{
		CaseID:      caseID,
		Source:      RecommendationSourceDoctor,
		ServiceCode: service.Code,
		ServiceName: service.Name,
		Importance:  importance,
		Rationale:   rationale,
		PatientText: patientText,
		Review:      review,
		CreatedAt:   review.ReviewedAt,
	}, nil
}

func (r Recommendation) Reviewed() bool {
	return r.Review.Mark != ""
}

func (r *Recommendation) Edit(review *Review, patientText *string) error {
	if review == nil && patientText == nil {
		return apperr.Invalid("mark or patient_text is required")
	}
	if review != nil {
		if err := review.Validate(); err != nil {
			return err
		}
	}
	if patientText != nil && strings.TrimSpace(*patientText) == "" {
		return apperr.Invalid("patient_text must not be empty")
	}

	if review != nil {
		r.Review = *review
	}
	if patientText != nil {
		r.PatientText = *patientText
	}
	return nil
}

func (r Recommendation) Accepted() bool {
	return r.Reviewed() && r.Review.Mark != MarkRejected
}

type Service struct {
	Code string
	Name string
}
