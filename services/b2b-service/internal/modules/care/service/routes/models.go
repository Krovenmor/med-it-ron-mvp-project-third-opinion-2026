package routes

import (
	careapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	doctorapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/api"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
)

func patientOf(e gatewayapi.PatientRegistered) domain.Patient {
	return domain.Patient{
		ID:           e.PatientID,
		SourceSystem: e.SourceSystem,
		ExternalID:   e.ExternalID,
		FullName:     e.FullName,
		BirthDate:    e.BirthDate,
		Sex:          e.Sex,
		Phone:        e.Phone,
		Email:        e.Email,
	}
}

func reviewStartOf(e doctorapi.ReviewStarted) domain.ReviewStart {
	return domain.ReviewStart{
		CaseID:      e.CaseID,
		PatientID:   e.PatientID,
		Urgency:     domain.Urgency(e.Urgency),
		Modality:    domain.Modality(e.Modality),
		PerformedAt: e.PerformedAt,
		ReceivedAt:  e.ReceivedAt,
		AssessedAt:  e.AssessedAt,
		ReviewDueAt: e.ReviewDueAt,
	}
}

func planItemsOf(e doctorapi.CaseConfirmed) []domain.PlanItem {
	items := make([]domain.PlanItem, 0, len(e.Recommendations))
	for _, r := range e.Recommendations {
		items = append(items, domain.PlanItem{
			ID:          r.ID,
			CaseID:      e.CaseID,
			Position:    r.Position,
			ServiceCode: r.ServiceCode,
			ServiceName: r.ServiceName,
			PatientText: r.PatientText,
			Mark:        domain.Mark(r.Mark),
		})
	}
	return items
}

func servicesOf(services []domain.Service) []careapi.Service {
	out := make([]careapi.Service, 0, len(services))
	for _, s := range services {
		out = append(out, careapi.Service{Code: s.Code, Name: s.Name})
	}
	return out
}
