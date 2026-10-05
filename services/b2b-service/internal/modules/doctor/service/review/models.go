package review

import (
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/api"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
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

func urgencyChanged(c domain.Case, from domain.Urgency, cmd ChangeUrgency) api.UrgencyChanged {
	return api.UrgencyChanged{
		CaseID:      c.ID,
		From:        string(from),
		To:          string(c.Urgency),
		Reason:      cmd.Reason,
		Actor:       cmd.Actor,
		ReviewDueAt: c.ReviewDueAt(),
		ChangedAt:   c.UpdatedAt,
	}
}

func caseConfirmed(c domain.Case, recs []domain.Recommendation, actor string) api.CaseConfirmed {
	accepted := make([]api.ConfirmedRecommendation, 0, len(recs))
	for _, r := range recs {
		if !r.Accepted() {
			continue
		}
		accepted = append(accepted, api.ConfirmedRecommendation{
			ID:          r.ID,
			Position:    r.Position,
			ServiceCode: r.ServiceCode,
			ServiceName: r.ServiceName,
			PatientText: r.PatientText,
			Mark:        string(r.Review.Mark),
		})
	}
	return api.CaseConfirmed{
		CaseID:          c.ID,
		Urgency:         string(c.Urgency),
		Actor:           actor,
		ConfirmedAt:     c.ConfirmedAt,
		Recommendations: accepted,
	}
}
