package intake

import (
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
)

const KindAssess = "assess"

type IngestResult struct {
	Intake  domain.Intake
	Created bool
}

func patientRegistered(p domain.Patient) api.PatientRegistered {
	return api.PatientRegistered{
		PatientID:    p.ID,
		SourceSystem: p.SourceSystem,
		ExternalID:   p.ExternalID,
		FullName:     p.FullName,
		BirthDate:    p.BirthDate,
		Sex:          string(p.Sex),
		Phone:        p.Phone,
		Email:        p.Email,
	}
}

func reportReceived(i domain.Intake) api.ReportReceived {
	return api.ReportReceived{
		CaseID:      i.CaseID,
		PatientID:   i.PatientID,
		Modality:    string(i.Study.Modality),
		PerformedAt: i.Study.PerformedAt,
		ReceivedAt:  i.ReceivedAt,
	}
}

func caseAssessed(i domain.Intake, p domain.Patient, a domain.Assessment) api.CaseAssessed {
	recs := make([]api.AssessedRecommendation, 0, len(a.Recommendations))
	for _, r := range a.Recommendations {
		recs = append(recs, api.AssessedRecommendation{
			ServiceCode:   r.ServiceCode,
			ServiceName:   r.ServiceName,
			Importance:    string(r.Importance),
			Rationale:     r.Rationale,
			GuidelineRef:  r.GuidelineRef,
			PatientText:   r.PatientText,
			AlreadyBooked: r.AlreadyBooked,
		})
	}
	return api.CaseAssessed{
		CaseID: i.CaseID,
		Patient: api.AssessedPatient{
			ID:           p.ID,
			SourceSystem: p.SourceSystem,
			ExternalID:   p.ExternalID,
			BirthDate:    p.BirthDate,
			Sex:          string(p.Sex),
		},
		Study: api.AssessedStudy{
			ID:          i.Study.ID,
			Modality:    string(i.Study.Modality),
			BodySite:    i.Study.BodySite,
			PerformedAt: i.Study.PerformedAt,
		},
		Conclusion:        i.Conclusion,
		ReceivedAt:        i.ReceivedAt,
		AssessedAt:        i.AssessedAt,
		Urgency:           string(a.Urgency),
		CatalogVersion:    a.CatalogVersion,
		GuidelinesVersion: a.GuidelinesVersion,
		Recommendations:   recs,
	}
}
