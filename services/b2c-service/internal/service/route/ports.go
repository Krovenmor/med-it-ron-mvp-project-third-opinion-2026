package route

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Plans interface {
	Plan(ctx context.Context, patientID string) (domain.Plan, error)
}

type MIS interface {
	History(ctx context.Context, patientID string) (domain.History, error)
}
