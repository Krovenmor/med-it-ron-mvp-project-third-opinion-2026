package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/history"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type caseRow struct {
	ID                  uuid.UUID  `db:"id"`
	PatientID           uuid.UUID  `db:"patient_id"`
	PatientSourceSystem string     `db:"patient_source_system"`
	PatientExternalID   string     `db:"patient_external_id"`
	PatientBirthDate    time.Time  `db:"patient_birth_date"`
	PatientSex          string     `db:"patient_sex"`
	StudyID             string     `db:"study_id"`
	Modality            string     `db:"modality"`
	BodySite            string     `db:"body_site"`
	PerformedAt         time.Time  `db:"performed_at"`
	Conclusion          string     `db:"conclusion"`
	Status              string     `db:"status"`
	Urgency             string     `db:"urgency"`
	CatalogVersion      string     `db:"catalog_version"`
	GuidelinesVersion   string     `db:"guidelines_version"`
	ReceivedAt          time.Time  `db:"received_at"`
	AssessedAt          time.Time  `db:"assessed_at"`
	ConfirmedAt         *time.Time `db:"confirmed_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}

type recommendationRow struct {
	ID            uuid.UUID  `db:"id"`
	CaseID        uuid.UUID  `db:"case_id"`
	Position      int        `db:"position"`
	Source        string     `db:"source"`
	ServiceCode   string     `db:"service_code"`
	ServiceName   string     `db:"service_name"`
	Importance    string     `db:"importance"`
	Rationale     string     `db:"rationale"`
	GuidelineRef  string     `db:"guideline_ref"`
	PatientText   string     `db:"patient_text"`
	AlreadyBooked bool       `db:"already_booked"`
	Mark          *string    `db:"mark"`
	RejectReason  *string    `db:"reject_reason"`
	RejectComment *string    `db:"reject_comment"`
	MarkedBy      *string    `db:"marked_by"`
	MarkedAt      *time.Time `db:"marked_at"`
	CreatedAt     time.Time  `db:"created_at"`
}

type reviewQueueRow struct {
	CaseID                  uuid.UUID  `db:"case_id"`
	Urgency                 string     `db:"urgency"`
	Modality                string     `db:"modality"`
	PerformedAt             time.Time  `db:"performed_at"`
	ReceivedAt              time.Time  `db:"received_at"`
	PatientID               uuid.UUID  `db:"patient_id"`
	PatientExternalID       string     `db:"patient_external_id"`
	PatientBirthDate        time.Time  `db:"patient_birth_date"`
	PatientSex              string     `db:"patient_sex"`
	RecommendationsTotal    int        `db:"recommendations_total"`
	RecommendationsReviewed int        `db:"recommendations_reviewed"`
	OpenedAt                *time.Time `db:"opened_at"`
}

func (r caseRow) toDomain() domain.Case {
	return domain.Case{
		ID: r.ID,
		Patient: domain.Patient{
			ID:           r.PatientID,
			SourceSystem: r.PatientSourceSystem,
			ExternalID:   r.PatientExternalID,
			BirthDate:    r.PatientBirthDate,
			Sex:          domain.Sex(r.PatientSex),
		},
		Study: domain.Study{
			ID:          r.StudyID,
			Modality:    domain.Modality(r.Modality),
			BodySite:    r.BodySite,
			PerformedAt: r.PerformedAt,
		},
		Conclusion:        r.Conclusion,
		Status:            domain.CaseStatus(r.Status),
		Urgency:           domain.Urgency(r.Urgency),
		CatalogVersion:    r.CatalogVersion,
		GuidelinesVersion: r.GuidelinesVersion,
		ReceivedAt:        r.ReceivedAt,
		AssessedAt:        r.AssessedAt,
		ConfirmedAt:       postgres.ValueOf(r.ConfirmedAt),
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
		Review: domain.Review{
			Mark:          domain.Mark(postgres.ValueOf(r.Mark)),
			RejectReason:  domain.RejectReason(postgres.ValueOf(r.RejectReason)),
			RejectComment: postgres.ValueOf(r.RejectComment),
			ReviewedBy:    postgres.ValueOf(r.MarkedBy),
			ReviewedAt:    postgres.ValueOf(r.MarkedAt),
		},
		CreatedAt: r.CreatedAt,
	}
}

func (r reviewQueueRow) toDomain() domain.ReviewQueueItem {
	return domain.ReviewQueueItem{
		CaseID:      r.CaseID,
		Urgency:     domain.Urgency(r.Urgency),
		Modality:    domain.Modality(r.Modality),
		PerformedAt: r.PerformedAt,
		ReceivedAt:  r.ReceivedAt,
		Patient: domain.Patient{
			ID:         r.PatientID,
			ExternalID: r.PatientExternalID,
			BirthDate:  r.PatientBirthDate,
			Sex:        domain.Sex(r.PatientSex),
		},
		RecommendationsTotal:    r.RecommendationsTotal,
		RecommendationsReviewed: r.RecommendationsReviewed,
		OpenedAt:                postgres.ValueOf(r.OpenedAt),
	}
}

func insertCaseArgs(c domain.Case) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":                    c.ID,
		"patient_id":            c.Patient.ID,
		"patient_source_system": c.Patient.SourceSystem,
		"patient_external_id":   c.Patient.ExternalID,
		"patient_birth_date":    c.Patient.BirthDate,
		"patient_sex":           string(c.Patient.Sex),
		"study_id":              c.Study.ID,
		"modality":              string(c.Study.Modality),
		"body_site":             c.Study.BodySite,
		"performed_at":          c.Study.PerformedAt,
		"conclusion":            c.Conclusion,
		"status":                string(c.Status),
		"urgency":               string(c.Urgency),
		"catalog_version":       c.CatalogVersion,
		"guidelines_version":    c.GuidelinesVersion,
		"received_at":           c.ReceivedAt,
		"assessed_at":           c.AssessedAt,
		"confirmed_at":          postgres.NullableTime(c.ConfirmedAt),
		"updated_at":            c.UpdatedAt,
	}
}

func updateCaseArgs(c domain.Case) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":           c.ID,
		"status":       string(c.Status),
		"urgency":      string(c.Urgency),
		"confirmed_at": postgres.NullableTime(c.ConfirmedAt),
		"updated_at":   c.UpdatedAt,
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

func appendRecommendationArgs(r domain.Recommendation) pgx.StrictNamedArgs {
	args := insertRecommendationArgs(r)
	delete(args, "position")
	args["mark"] = string(r.Review.Mark)
	args["marked_by"] = r.Review.ReviewedBy
	args["marked_at"] = r.Review.ReviewedAt
	return args
}

func updateReviewArgs(r domain.Recommendation) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":             r.ID,
		"mark":           string(r.Review.Mark),
		"reject_reason":  string(r.Review.RejectReason),
		"reject_comment": r.Review.RejectComment,
		"marked_by":      r.Review.ReviewedBy,
		"marked_at":      r.Review.ReviewedAt,
		"patient_text":   r.PatientText,
	}
}

func insertCaseEventArgs(e domain.CaseEvent) pgx.StrictNamedArgs {
	payload := e.Payload
	if payload == nil {
		payload = map[string]string{}
	}
	return pgx.StrictNamedArgs{
		"payload":     payload,
		"case_id":     e.CaseID,
		"type":        string(e.Type),
		"from_status": string(e.FromStatus),
		"to_status":   string(e.ToStatus),
		"actor":       e.Actor,
		"occurred_at": e.OccurredAt,
	}
}

func historyCaseArgs(c history.Case) pgx.StrictNamedArgs {
	return insertCaseArgs(domain.Case{
		ID: c.ID,
		Patient: domain.Patient{
			ID:           c.Patient.ID,
			SourceSystem: c.SourceSystem,
			ExternalID:   c.Patient.ExternalID,
			BirthDate:    c.Patient.BirthDate,
			Sex:          domain.Sex(c.Patient.Sex),
		},
		Study: domain.Study{
			ID:          c.Study.ID,
			Modality:    domain.Modality(c.Study.Modality),
			BodySite:    c.Study.BodySite,
			PerformedAt: c.Study.PerformedAt,
		},
		Conclusion:        c.Conclusion,
		Status:            domain.CaseStatusConfirmed,
		Urgency:           domain.Urgency(c.Urgency),
		CatalogVersion:    c.CatalogVersion,
		GuidelinesVersion: c.GuidelinesVersion,
		ReceivedAt:        c.ReceivedAt,
		AssessedAt:        c.AssessedAt,
		ConfirmedAt:       c.Review.ConfirmedAt,
		UpdatedAt:         c.Review.ConfirmedAt,
	})
}

func historyRecommendationArgs(c history.Case, r history.Recommendation) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":             r.ID,
		"case_id":        c.ID,
		"position":       r.Position,
		"source":         string(domain.RecommendationSourceAI),
		"service_code":   r.ServiceCode,
		"service_name":   r.ServiceName,
		"importance":     r.Importance,
		"rationale":      r.Rationale,
		"guideline_ref":  r.GuidelineRef,
		"patient_text":   r.PatientText,
		"already_booked": false,
		"mark":           r.Mark,
		"reject_reason":  r.RejectReason,
		"reject_comment": r.RejectComment,
		"marked_by":      c.Review.Doctor,
		"marked_at":      postgres.NullableTime(c.Review.OpenedAt.Add(c.Review.ConfirmedAt.Sub(c.Review.OpenedAt) / 2)),
		"created_at":     c.AssessedAt,
	}
}

func historyEvents(c history.Case) []domain.CaseEvent {
	return []domain.CaseEvent{
		domain.StatusChanged(c.ID, "", domain.CaseStatusInReview, domain.ActorSystem, c.AssessedAt),
		domain.CaseOpened(c.ID, c.Review.Doctor, c.Review.OpenedAt),
		domain.StatusChanged(c.ID, domain.CaseStatusInReview, domain.CaseStatusConfirmed, c.Review.Doctor, c.Review.ConfirmedAt),
	}
}
