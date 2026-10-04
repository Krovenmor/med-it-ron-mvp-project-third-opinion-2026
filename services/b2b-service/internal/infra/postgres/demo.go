package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres/queries"
)

type Demo struct {
	executor
	q queries.Demo
}

func NewDemo(pool *pgxpool.Pool, q queries.Queries) *Demo {
	return &Demo{executor: executor{pool}, q: q.Demo}
}

func (r *Demo) Reset(ctx context.Context) error {
	if _, err := r.db(ctx).Exec(ctx, r.q.Reset); err != nil {
		return fmt.Errorf("reset demo data: %w", err)
	}
	return nil
}
