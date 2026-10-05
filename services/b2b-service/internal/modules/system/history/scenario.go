package history

import (
	"time"

	"github.com/google/uuid"
)

const demoSourceSystem = "mis-demo"

var clinicZone = time.FixedZone("MSK", 3*60*60)

func DemoScenario(now time.Time) []Case {
	return []Case{igorPreviousCase(anchorDay(now))}
}

func igorPreviousCase(day time.Time) Case {
	at := func(days, hour int) time.Time {
		return day.AddDate(0, 0, days).Add(time.Duration(hour) * time.Hour).UTC()
	}
	performedAt := at(-97, 9)
	assessedAt := performedAt.Add(20 * time.Minute)
	confirmedAt := performedAt.Add(5 * time.Hour)
	pulmonology := Recommendation{
		ID:          uuid.New(),
		Position:    1,
		ServiceCode: "PULM-CONSULT",
		ServiceName: "Консультация пульмонолога",
		Importance:  "high",
		Rationale:   "Очаг в лёгком требует оценки специалистом.",
		PatientText: "Рекомендуем консультацию пульмонолога по результатам КТ.",
		Mark:        "critical",
	}
	followUp := Recommendation{
		ID:          uuid.New(),
		Position:    2,
		ServiceCode: "CT-CHEST-FOLLOWUP",
		ServiceName: "Контрольная КТ органов грудной клетки",
		Importance:  "low",
		Rationale:   "Динамическое наблюдение очага.",
		PatientText: "Повторное исследование покажет, как изменилась картина со временем.",
		Mark:        "minor",
	}
	return Case{
		ID:           uuid.New(),
		SourceSystem: demoSourceSystem,
		Patient: Patient{
			ID:         uuid.New(),
			ExternalID: "P-IGOR",
			FullName:   "Игорь Петров",
			BirthDate:  time.Date(1965, 5, 1, 0, 0, 0, 0, time.UTC),
			Sex:        "male",
			Phone:      "+79990000002",
		},
		Study: Study{
			ID:          "CT-2026-HIST-0001",
			Modality:    "CT",
			BodySite:    "chest",
			PerformedAt: performedAt,
		},
		Conclusion:        "КТ органов грудной клетки. В S6 правого лёгкого очаг 6 мм с ровными контурами.",
		Urgency:           "planned",
		CatalogVersion:    catalogVersion,
		GuidelinesVersion: guidelinesVersion,
		ReceivedAt:        performedAt.Add(10 * time.Minute),
		AssessedAt:        assessedAt,
		ReviewDueAt:       assessedAt.Add(reviewSLA("planned")),
		Review: Review{
			Doctor:      "dr-demo",
			OpenedAt:    performedAt.Add(2 * time.Hour),
			ConfirmedAt: confirmedAt,
		},
		Recommendations: []Recommendation{pulmonology, followUp},
		NotifiedAt:      confirmedAt.Add(10 * time.Second),
		Bookings: []Booking{{
			ID:               uuid.New(),
			RecommendationID: pulmonology.ID,
			AppointmentID:    "HIST-APT-" + pulmonology.ID.String(),
			ScheduledAt:      at(-85, 14),
			Channel:          "self",
			CreatedAt:        at(-96, 12),
		}},
	}
}

func anchorDay(now time.Time) time.Time {
	local := now.In(clinicZone)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, clinicZone)
}
