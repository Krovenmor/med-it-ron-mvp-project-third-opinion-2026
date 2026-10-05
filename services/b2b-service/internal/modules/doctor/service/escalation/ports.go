package escalation

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
)

type Clock interface {
	Now() time.Time
}

type Cases interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Case, error)
}

type Events interface {
	Exists(ctx context.Context, caseID uuid.UUID, eventType domain.CaseEventType) (bool, error)
}

type Publisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
}
