package aiservice

import (
	"context"
	"strings"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

const mockVersion = "mock-2026.10"

type Mock struct{}

func NewMock() Mock {
	return Mock{}
}

func (Mock) Assess(_ context.Context, req domain.AssessmentRequest) (domain.Assessment, error) {
	conclusion := strings.ToLower(req.Conclusion)
	switch {
	case strings.Contains(conclusion, "пневмоторакс"):
		return mockAssessment(domain.UrgencyEmergency, domain.Recommendation{
			ServiceCode:  "SURG-THOR-CONSULT",
			ServiceName:  "Консультация торакального хирурга",
			Importance:   domain.ImportanceHigh,
			Rationale:    "Признаки пневмоторакса требуют неотложной оценки хирургом.",
			GuidelineRef: "КР МЗ РФ «Пневмоторакс»",
			PatientText:  "Врач просит вас срочно связаться с клиникой.",
		}), nil
	case req.Modality == domain.ModalityMG && strings.Contains(conclusion, "bi-rads 4"):
		return mockAssessment(domain.UrgencyPriority,
			domain.Recommendation{
				ServiceCode:  "ONC-MAMMO-CONSULT",
				ServiceName:  "Консультация онколога-маммолога",
				Importance:   domain.ImportanceHigh,
				Rationale:    "BI-RADS 4 — показана консультация специалиста и решение вопроса о биопсии.",
				GuidelineRef: "КР МЗ РФ «Рак молочной железы»",
				PatientText:  "Рекомендуем консультацию маммолога, чтобы обсудить результаты обследования и дальнейшие шаги.",
			},
			domain.Recommendation{
				ServiceCode:  "US-BREAST",
				ServiceName:  "УЗИ молочных желез",
				Importance:   domain.ImportanceLow,
				Rationale:    "Дополнительная визуализация для уточнения характеристик образования.",
				GuidelineRef: "КР МЗ РФ «Рак молочной железы»",
				PatientText:  "Ультразвуковое исследование поможет врачу получить больше информации.",
			},
		), nil
	case strings.Contains(conclusion, "очаг"):
		return mockAssessment(domain.UrgencyPriority,
			domain.Recommendation{
				ServiceCode:  "PULM-CONSULT",
				ServiceName:  "Консультация пульмонолога",
				Importance:   domain.ImportanceHigh,
				Rationale:    "Солидный очаг более 8 мм требует оценки специалистом и определения тактики наблюдения.",
				GuidelineRef: "КР МЗ РФ «Злокачественное новообразование бронхов и легкого»",
				PatientText:  "Рекомендуем консультацию пульмонолога по результатам КТ.",
			},
			domain.Recommendation{
				ServiceCode:  "CT-CHEST-FOLLOWUP",
				ServiceName:  "Контрольная КТ органов грудной клетки через 3 месяца",
				Importance:   domain.ImportanceLow,
				Rationale:    "Динамическое наблюдение очага.",
				GuidelineRef: "КР МЗ РФ «Злокачественное новообразование бронхов и легкого»",
				PatientText:  "Повторное исследование покажет, как изменилась картина со временем.",
			},
		), nil
	default:
		return mockAssessment(domain.UrgencyNormal, domain.Recommendation{
			ServiceCode: "SCREENING-ANNUAL",
			ServiceName: "Плановый профилактический осмотр",
			Importance:  domain.ImportanceLow,
			Rationale:   "Патологии не выявлено, плановый скрининг в установленный срок.",
			PatientText: "Всё в порядке. Напомним вам о следующем плановом обследовании.",
		}), nil
	}
}

func mockAssessment(urgency domain.Urgency, recs ...domain.Recommendation) domain.Assessment {
	return domain.Assessment{
		Urgency:           urgency,
		CatalogVersion:    mockVersion,
		GuidelinesVersion: mockVersion,
		Recommendations:   recs,
	}
}
