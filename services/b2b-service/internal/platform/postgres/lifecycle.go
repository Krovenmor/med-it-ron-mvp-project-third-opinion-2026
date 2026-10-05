package postgres

import (
	"context"
	"io/fs"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	DSN            string        `env:"DSN"`
	ConnectTimeout time.Duration `env:"CONNECT_TIMEOUT"`
}

type Schema struct {
	Name       string
	Create     string
	Migrations fs.FS
}

func OpenModule(lc fx.Lifecycle, log *zap.Logger, cfg Config, schema Schema) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()
	db, err := Open(ctx, cfg.DSN, schema.Name)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			applied, err := db.Migrate(ctx, schema.Create, schema.Migrations)
			if err != nil {
				return err
			}
			log.Info("database migrated", zap.String("schema", schema.Name), zap.Strings("applied", applied))
			return nil
		},
		OnStop: func(context.Context) error {
			db.Pool.Close()
			return nil
		},
	})
	return db, nil
}
