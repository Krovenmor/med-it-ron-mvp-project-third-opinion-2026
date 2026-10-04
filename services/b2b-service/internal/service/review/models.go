package review

import (
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type ReviewRecommendation struct {
	CaseID           uuid.UUID
	RecommendationID uuid.UUID
	Mark             domain.Mark
	RejectReason     domain.RejectReason
	RejectComment    string
	PatientText      *string
	Actor            string
}

type AddRecommendation struct {
	CaseID      uuid.UUID
	Service     domain.Service
	Rationale   string
	PatientText string
	Mark        domain.Mark
	Actor       string
}

type ChangeUrgency struct {
	CaseID  uuid.UUID
	Urgency domain.Urgency
	Reason  string
	Actor   string
}
