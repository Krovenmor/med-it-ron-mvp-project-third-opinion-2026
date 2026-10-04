package migrations

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var files embed.FS

func Up(ctx context.Context, pool *pgxpool.Pool) ([]string, error) {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, files)
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
