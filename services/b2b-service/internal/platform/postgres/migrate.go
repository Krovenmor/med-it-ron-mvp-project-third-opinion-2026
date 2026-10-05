package postgres

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (d *Database) Migrate(ctx context.Context, createSchema string, migrations fs.FS) ([]string, error) {
	if _, err := d.Pool.Exec(ctx, createSchema); err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}

	db := stdlib.OpenDBFromPool(d.Pool)
	defer db.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations)
	if err != nil {
		return nil, fmt.Errorf("create migration provider: %w", err)
	}
	results, err := provider.Up(ctx)
	if err != nil {
		return nil, fmt.Errorf("apply migrations: %w", err)
	}

	applied := make([]string, 0, len(results))
	for _, r := range results {
		applied = append(applied, r.Source.Path)
	}
	return applied, nil
}
