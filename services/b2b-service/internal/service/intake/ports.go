package intake

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Clock interface {
	Now() time.Time
}

type Patients interface {
	Upsert(ctx context.Context, p domain.Patient, now time.Time) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Patient, error)
}

type Cases interface {
	CreateIfAbsent(ctx context.Context, c domain.Case) (uuid.UUID, bool, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Case, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (domain.Case, error)
	GetByStudy(ctx context.Context, sourceSystem, studyID string) (domain.Case, error)
	Update(ctx context.Context, c domain.Case) error
}

type Recommendations interface {
	InsertMany(ctx context.Context, recs []domain.Recommendation) error
	ActiveServices(ctx context.Context, patientID, excludeCaseID uuid.UUID) ([]domain.Service, error)
}

type Events interface {
	Append(ctx context.Context, e domain.CaseEvent) error
}

type Jobs interface {
	Enqueue(ctx context.Context, job domain.Job, now time.Time) error
}

type AIService interface {
	Assess(ctx context.Context, req domain.AssessmentRequest) (domain.Assessment, error)
}

type MIS interface {
	PatientHistory(ctx context.Context, p domain.Patient) (domain.PatientHistory, error)
}
