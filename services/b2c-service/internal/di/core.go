package di

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/config"
)

var coreModule = fx.Module("core",
	fx.Provide(
		config.Load,
		newLogger,
	),
)

func newLogger(cfg config.Config) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	zcfg := zap.NewProductionConfig()
	zcfg.Level = zap.NewAtomicLevelAt(level)
	return zcfg.Build()
}

func newFxLogger(log *zap.Logger) fxevent.Logger {
	l := &fxevent.ZapLogger{Logger: log.Named("fx")}
	l.UseLogLevel(zapcore.DebugLevel)
	return l
}
