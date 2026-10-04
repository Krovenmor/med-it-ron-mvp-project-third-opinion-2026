package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
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

type caseRow struct {
	ID                uuid.UUID `db:"id"`
	PatientID         uuid.UUID `db:"patient_id"`
	SourceSystem      string    `db:"source_system"`
	StudyID           string    `db:"study_id"`
	Modality          string    `db:"modality"`
	BodySite          string    `db:"body_site"`
	PerformedAt       time.Time `db:"performed_at"`
	Conclusion        string    `db:"conclusion"`
	Fingerprint       []byte    `db:"fingerprint"`
	Status            string    `db:"status"`
	Urgency           *string   `db:"urgency"`
	CatalogVersion    string    `db:"catalog_version"`
	GuidelinesVersion string    `db:"guidelines_version"`
	ReceivedAt        time.Time `db:"received_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

type recommendationRow struct {
	ID            uuid.UUID `db:"id"`
	CaseID        uuid.UUID `db:"case_id"`
	Position      int       `db:"position"`
	Source        string    `db:"source"`
	ServiceCode   string    `db:"service_code"`
	ServiceName   string    `db:"service_name"`
	Importance    string    `db:"importance"`
	Rationale     string    `db:"rationale"`
	GuidelineRef  string    `db:"guideline_ref"`
	PatientText   string    `db:"patient_text"`
	AlreadyBooked bool      `db:"already_booked"`
	CreatedAt     time.Time `db:"created_at"`
}

type serviceRow struct {
	Code string `db:"service_code"`
	Name string `db:"service_name"`
}

type jobRow struct {
	ID       int64     `db:"id"`
	CaseID   uuid.UUID `db:"case_id"`
	Kind     string    `db:"kind"`
	RunAt    time.Time `db:"run_at"`
	Attempts int       `db:"attempts"`
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

func (r caseRow) toDomain() domain.Case {
	var urgency domain.Urgency
	if r.Urgency != nil {
		urgency = domain.Urgency(*r.Urgency)
	}
	return domain.Case{
		ID:           r.ID,
		PatientID:    r.PatientID,
		SourceSystem: r.SourceSystem,
		Study: domain.Study{
			ID:          r.StudyID,
			Modality:    domain.Modality(r.Modality),
			BodySite:    r.BodySite,
			PerformedAt: r.PerformedAt,
		},
		Conclusion:        r.Conclusion,
		Fingerprint:       r.Fingerprint,
		Status:            domain.CaseStatus(r.Status),
		Urgency:           urgency,
		CatalogVersion:    r.CatalogVersion,
		GuidelinesVersion: r.GuidelinesVersion,
		ReceivedAt:        r.ReceivedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

func (r recommendationRow) toDomain() domain.Recommendation {
	return domain.Recommendation{
		ID:            r.ID,
		CaseID:        r.CaseID,
		Position:      r.Position,
		Source:        domain.RecommendationSource(r.Source),
		ServiceCode:   r.ServiceCode,
		ServiceName:   r.ServiceName,
		Importance:    domain.Importance(r.Importance),
		Rationale:     r.Rationale,
		GuidelineRef:  r.GuidelineRef,
		PatientText:   r.PatientText,
		AlreadyBooked: r.AlreadyBooked,
		CreatedAt:     r.CreatedAt,
	}
}

func (r serviceRow) toDomain() domain.Service {
	return domain.Service{Code: r.Code, Name: r.Name}
}

func (r jobRow) toDomain() domain.Job {
	return domain.Job{
		ID:       r.ID,
		CaseID:   r.CaseID,
		Kind:     domain.JobKind(r.Kind),
		RunAt:    r.RunAt,
		Attempts: r.Attempts,
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

func insertCaseArgs(c domain.Case) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"patient_id":    c.PatientID,
		"source_system": c.SourceSystem,
		"study_id":      c.Study.ID,
		"modality":      string(c.Study.Modality),
		"body_site":     c.Study.BodySite,
		"performed_at":  c.Study.PerformedAt,
		"conclusion":    c.Conclusion,
		"fingerprint":   c.Fingerprint,
		"status":        string(c.Status),
		"received_at":   c.ReceivedAt,
		"updated_at":    c.UpdatedAt,
	}
}

func updateCaseArgs(c domain.Case) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":                 c.ID,
		"status":             string(c.Status),
		"urgency":            string(c.Urgency),
		"catalog_version":    c.CatalogVersion,
		"guidelines_version": c.GuidelinesVersion,
		"updated_at":         c.UpdatedAt,
	}
}

func insertRecommendationArgs(r domain.Recommendation) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"case_id":        r.CaseID,
		"position":       r.Position,
		"source":         string(r.Source),
		"service_code":   r.ServiceCode,
		"service_name":   r.ServiceName,
		"importance":     string(r.Importance),
		"rationale":      r.Rationale,
		"guideline_ref":  r.GuidelineRef,
		"patient_text":   r.PatientText,
		"already_booked": r.AlreadyBooked,
		"created_at":     r.CreatedAt,
	}
}

func insertCaseEventArgs(e domain.CaseEvent) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"case_id":     e.CaseID,
		"type":        string(e.Type),
		"from_status": string(e.FromStatus),
		"to_status":   string(e.ToStatus),
		"actor":       e.Actor,
		"occurred_at": e.OccurredAt,
	}
}

func enqueueJobArgs(job domain.Job, now time.Time) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"case_id":    job.CaseID,
		"kind":       string(job.Kind),
		"run_at":     job.RunAt,
		"created_at": now,
	}
}

func claimJobArgs(now time.Time, lease time.Duration) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"now":           now,
		"lease_seconds": lease.Seconds(),
	}
}

func completeJobArgs(job domain.Job) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":       job.ID,
		"attempts": job.Attempts,
	}
}

func rescheduleJobArgs(job domain.Job, runAt time.Time, cause string) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":         job.ID,
		"attempts":   job.Attempts,
		"run_at":     runAt,
		"last_error": cause,
	}
}

func failJobArgs(job domain.Job, cause string) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":         job.ID,
		"attempts":   job.Attempts,
		"last_error": cause,
	}
}
