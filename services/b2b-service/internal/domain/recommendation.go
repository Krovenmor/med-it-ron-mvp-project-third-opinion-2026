package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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
			return invalid("reject_reason and reject_comment are allowed only for rejected recommendations")
		}
	case MarkRejected:
		switch r.RejectReason {
		case RejectReasonContraindicated:
		case RejectReasonOther:
			if r.RejectComment == "" {
				return invalid("reject_comment is required when reject_reason is other")
			}
		default:
			return invalid("reject_reason must be one of contraindicated, other")
		}
	default:
		return invalid("mark must be one of critical, minor, rejected")
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
		return Recommendation{}, invalid("service_name is required")
	case patientText == "":
		return Recommendation{}, invalid("patient_text is required")
	case review.Mark != MarkCritical && review.Mark != MarkMinor:
		return Recommendation{}, invalid("mark must be one of critical, minor")
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

func (r *Recommendation) ApplyReview(review Review, patientText *string) error {
	if err := review.Validate(); err != nil {
		return err
	}
	if patientText != nil {
		if *patientText == "" {
			return invalid("patient_text must not be empty")
		}
		r.PatientText = *patientText
	}
	r.Review = review
	return nil
}

func (r Recommendation) VisibleToPatient() bool {
	return r.Reviewed() && r.Review.Mark != MarkRejected
}

func (r Recommendation) EnsureBookable() error {
	switch {
	case !r.VisibleToPatient():
		return fmt.Errorf("%w: recommendation is not confirmed by a doctor", ErrInvalidState)
	case r.ServiceCode == "":
		return fmt.Errorf("%w: service is not available in the clinic", ErrInvalidState)
	}
	return nil
}

type Service struct {
	Code string
	Name string
}
