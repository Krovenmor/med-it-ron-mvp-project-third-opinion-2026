package postgres

import (
	"context"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct {
	schema string
}

type Database struct {
	Pool   *pgxpool.Pool
	Tx     trm.Manager
	getter *trmpgx.CtxGetter
}

func Open(ctx context.Context, dsn, schema string) (*Database, error) {
	pool, err := Connect(ctx, dsn, schema)
	if err != nil {
		return nil, err
	}
	key := txKey{schema: schema}
	ctxManager := trmcontext.New(key)
	tx, err := manager.New(
		trmpgx.NewDefaultFactory(pool),
		manager.WithSettings(trmpgx.MustSettings(settings.Must(settings.WithCtxKey(key)))),
		manager.WithCtxManager(ctxManager),
	)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create %s transaction manager: %w", schema, err)
	}
	return &Database{Pool: pool, Tx: tx, getter: trmpgx.NewCtxGetter(ctxManager)}, nil
}

func (d *Database) Executor() Executor {
	return Executor{pool: d.Pool, getter: d.getter}
}

func Connect(ctx context.Context, dsn, schema string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	cfg.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "timestamptz",
			OID:   pgtype.TimestamptzOID,
			Codec: &pgtype.TimestamptzCodec{ScanLocation: time.UTC},
		})
		return nil
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}

type Executor struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (e Executor) DB(ctx context.Context) trmpgx.Tr {
	return e.getter.DefaultTrOrDB(ctx, e.pool)
}

func ValueOf[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}

func NullableTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
