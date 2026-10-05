package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/history"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type patientRow struct {
	ID           uuid.UUID `db:"id"`
	SourceSystem string    `db:"source_system"`
	ExternalID   string    `db:"external_id"`
	FullName     string    `db:"full_name"`
	BirthDate    time.Time `db:"birth_date"`
	Sex          string    `db:"sex"`
	Phone        string    `db:"phone"`
	Email        string    `db:"email"`
}

type intakeRow struct {
	CaseID       uuid.UUID  `db:"case_id"`
	PatientID    uuid.UUID  `db:"patient_id"`
	SourceSystem string     `db:"source_system"`
	StudyID      string     `db:"study_id"`
	Modality     string     `db:"modality"`
	BodySite     string     `db:"body_site"`
	PerformedAt  time.Time  `db:"performed_at"`
	Conclusion   string     `db:"conclusion"`
	Fingerprint  []byte     `db:"fingerprint"`
	Status       string     `db:"status"`
	ReceivedAt   time.Time  `db:"received_at"`
	AssessedAt   *time.Time `db:"assessed_at"`
}

func (r patientRow) toDomain() domain.Patient {
	return domain.Patient{
		ID:           r.ID,
		SourceSystem: r.SourceSystem,
		ExternalID:   r.ExternalID,
		FullName:     r.FullName,
		BirthDate:    r.BirthDate,
		Sex:          domain.Sex(r.Sex),
		Phone:        r.Phone,
		Email:        r.Email,
	}
}

func (r intakeRow) toDomain() domain.Intake {
	return domain.Intake{
		CaseID:       r.CaseID,
		PatientID:    r.PatientID,
		SourceSystem: r.SourceSystem,
		Study: domain.Study{
			ID:          r.StudyID,
			Modality:    domain.Modality(r.Modality),
			BodySite:    r.BodySite,
			PerformedAt: r.PerformedAt,
		},
		Conclusion:  r.Conclusion,
		Fingerprint: r.Fingerprint,
		Status:      domain.IntakeStatus(r.Status),
		ReceivedAt:  r.ReceivedAt,
		AssessedAt:  postgres.ValueOf(r.AssessedAt),
	}
}

func upsertPatientArgs(p domain.Patient, now time.Time) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"source_system": p.SourceSystem,
		"external_id":   p.ExternalID,
		"full_name":     p.FullName,
		"birth_date":    p.BirthDate,
		"sex":           string(p.Sex),
		"phone":         p.Phone,
		"email":         p.Email,
		"now":           now,
	}
}

func insertIntakeArgs(i domain.Intake) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"patient_id":    i.PatientID,
		"source_system": i.SourceSystem,
		"study_id":      i.Study.ID,
		"modality":      string(i.Study.Modality),
		"body_site":     i.Study.BodySite,
		"performed_at":  i.Study.PerformedAt,
		"conclusion":    i.Conclusion,
		"fingerprint":   i.Fingerprint,
		"status":        string(i.Status),
		"received_at":   i.ReceivedAt,
	}
}

func updateIntakeArgs(i domain.Intake) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"case_id":     i.CaseID,
		"status":      string(i.Status),
		"assessed_at": postgres.NullableTime(i.AssessedAt),
	}
}

func historyPatientArgs(c history.Case) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":            c.Patient.ID,
		"source_system": c.SourceSystem,
		"external_id":   c.Patient.ExternalID,
		"full_name":     c.Patient.FullName,
		"birth_date":    c.Patient.BirthDate,
		"sex":           c.Patient.Sex,
		"phone":         c.Patient.Phone,
		"email":         "",
		"created_at":    c.ReceivedAt,
	}
}

func historyIntakeArgs(c history.Case) pgx.StrictNamedArgs {
	report := domain.Report{
		SourceSystem: c.SourceSystem,
		Study: domain.Study{
			ID:          c.Study.ID,
			Modality:    domain.Modality(c.Study.Modality),
			BodySite:    c.Study.BodySite,
			PerformedAt: c.Study.PerformedAt,
		},
		Conclusion: c.Conclusion,
		Patient:    domain.Patient{ExternalID: c.Patient.ExternalID},
	}
	return pgx.StrictNamedArgs{
		"case_id":       c.ID,
		"patient_id":    c.Patient.ID,
		"source_system": c.SourceSystem,
		"study_id":      c.Study.ID,
		"modality":      c.Study.Modality,
		"body_site":     c.Study.BodySite,
		"performed_at":  c.Study.PerformedAt,
		"conclusion":    c.Conclusion,
		"fingerprint":   report.Fingerprint(),
		"status":        string(domain.IntakeStatusAssessed),
		"received_at":   c.ReceivedAt,
		"assessed_at":   c.AssessedAt,
	}
}
