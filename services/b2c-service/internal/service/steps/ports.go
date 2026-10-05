package steps

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Plans interface {
	Plan(ctx context.Context, patientID string) (domain.Plan, error)
	DeclineRecommendation(ctx context.Context, caseID string, d domain.Decline) error
	RequestHelp(ctx context.Context, caseID string) error
}
