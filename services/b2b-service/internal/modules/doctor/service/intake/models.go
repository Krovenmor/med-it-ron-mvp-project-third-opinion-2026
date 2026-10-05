package intake

import (
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
)

func caseOf(e gatewayapi.CaseAssessed) domain.Case {
	return domain.Case{
		ID: e.CaseID,
		Patient: domain.Patient{
			ID:           e.Patient.ID,
			SourceSystem: e.Patient.SourceSystem,
			ExternalID:   e.Patient.ExternalID,
			BirthDate:    e.Patient.BirthDate,
			Sex:          domain.Sex(e.Patient.Sex),
		},
		Study: domain.Study{
			ID:          e.Study.ID,
			Modality:    domain.Modality(e.Study.Modality),
			BodySite:    e.Study.BodySite,
			PerformedAt: e.Study.PerformedAt,
		},
		Conclusion:        e.Conclusion,
		Status:            domain.CaseStatusInReview,
		Urgency:           domain.Urgency(e.Urgency),
		CatalogVersion:    e.CatalogVersion,
		GuidelinesVersion: e.GuidelinesVersion,
		ReceivedAt:        e.ReceivedAt,
		AssessedAt:        e.AssessedAt,
		UpdatedAt:         e.AssessedAt,
	}
}

func recommendationsOf(caseID uuid.UUID, e gatewayapi.CaseAssessed) []domain.Recommendation {
	recs := make([]domain.Recommendation, 0, len(e.Recommendations))
	for i, r := range e.Recommendations {
		recs = append(recs, domain.Recommendation{
			CaseID:        caseID,
			Position:      i + 1,
			Source:        domain.RecommendationSourceAI,
			ServiceCode:   r.ServiceCode,
			ServiceName:   r.ServiceName,
			Importance:    domain.Importance(r.Importance),
			Rationale:     r.Rationale,
			GuidelineRef:  r.GuidelineRef,
			PatientText:   r.PatientText,
			AlreadyBooked: r.AlreadyBooked,
			CreatedAt:     e.AssessedAt,
		})
	}
	return recs
}

func reviewStarted(c domain.Case) api.ReviewStarted {
	return api.ReviewStarted{
		CaseID:      c.ID,
		PatientID:   c.Patient.ID,
		Urgency:     string(c.Urgency),
		Modality:    string(c.Study.Modality),
		PerformedAt: c.Study.PerformedAt,
		ReceivedAt:  c.ReceivedAt,
		AssessedAt:  c.AssessedAt,
		ReviewDueAt: c.ReviewDueAt(),
	}
}
