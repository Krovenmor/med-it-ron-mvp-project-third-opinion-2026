package intake

import (
	"context"
	"time"

	"github.com/google/uuid"

	careapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

type Clock interface {
	Now() time.Time
}

type Patients interface {
	Upsert(ctx context.Context, p domain.Patient, now time.Time) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Patient, error)
}

type Intakes interface {
	CreateIfAbsent(ctx context.Context, i domain.Intake) (uuid.UUID, bool, error)
	Get(ctx context.Context, caseID uuid.UUID) (domain.Intake, error)
	GetForUpdate(ctx context.Context, caseID uuid.UUID) (domain.Intake, error)
	GetByStudy(ctx context.Context, sourceSystem, studyID string) (domain.Intake, error)
	Update(ctx context.Context, i domain.Intake) error
}

type Jobs interface {
	Enqueue(ctx context.Context, job jobs.Job, now time.Time) error
}

type Publisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
}

type AIService interface {
	Assess(ctx context.Context, req domain.AssessmentRequest) (domain.Assessment, error)
}

type MIS interface {
	PatientHistory(ctx context.Context, externalID string) (domain.PatientHistory, error)
}

type Routes interface {
	ActiveServices(ctx context.Context, patientID, excludeCaseID uuid.UUID) ([]careapi.Service, error)
}
