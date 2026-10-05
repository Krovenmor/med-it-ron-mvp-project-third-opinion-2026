package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Notifications struct {
	postgres.Executor
	q queries.Notifications
}

func NewNotifications(db *postgres.Database, q queries.Queries) *Notifications {
	return &Notifications{Executor: db.Executor(), q: q.Notifications}
}

func (r *Notifications) InsertMany(ctx context.Context, notifications []domain.Notification) error {
	batch := &pgx.Batch{}
	for _, n := range notifications {
		batch.Queue(r.q.Insert, insertNotificationArgs(n))
	}
	if err := r.DB(ctx).SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("insert notifications: %w", err)
	}
	return nil
}

func (r *Notifications) List(ctx context.Context, limit int) ([]domain.NotificationEntry, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.List, pgx.StrictNamedArgs{"limit": limit})
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	entries, err := postgres.CollectAll(rows, notificationEntryRow.toDomain)
	if err != nil {
		return nil, fmt.Errorf("scan notifications: %w", err)
	}
	return entries, nil
}
