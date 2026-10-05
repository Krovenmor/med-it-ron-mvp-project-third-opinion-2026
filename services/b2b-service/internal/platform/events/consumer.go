package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/events/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Inbox struct {
	postgres.Executor
	q queries.Queries
}

func NewInbox(db *postgres.Database) (*Inbox, error) {
	q, err := queries.Load()
	if err != nil {
		return nil, err
	}
	return &Inbox{Executor: db.Executor(), q: q}, nil
}

func (r *Inbox) Remember(ctx context.Context, m Message, now time.Time) (bool, error) {
	var id string
	err := r.DB(ctx).QueryRow(ctx, r.q.Remember, rememberArgs(m, now)).Scan(&id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("remember %s message: %w", m.Topic, err)
	}
	return true, nil
}

type Consumer struct {
	tx    trm.Manager
	inbox *Inbox
	clock Clock
}

func NewConsumer(tx trm.Manager, inbox *Inbox, clock Clock) *Consumer {
	return &Consumer{tx: tx, inbox: inbox, clock: clock}
}

func (c *Consumer) Subscribe(topic string, handle Handler) Subscription {
	return Subscription{Topic: topic, Handle: func(ctx context.Context, m Message) error {
		return c.tx.Do(ctx, func(ctx context.Context) error {
			fresh, err := c.inbox.Remember(ctx, m, c.clock.Now())
			if err != nil || !fresh {
				return err
			}
			return handle(ctx, m)
		})
	}}
}
