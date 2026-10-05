package demo

import (
	"context"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/history"
)

type Clock interface {
	Now() time.Time
	Reset()
}

type Module interface {
	Reset(ctx context.Context, cases []history.Case) error
}
