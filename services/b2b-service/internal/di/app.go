package di

import (
	"context"
	"errors"
	"net"
	"net/http"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/aiservice"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/mis"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres/migrations"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/cases"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/intake"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/review"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/http"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/worker"
)

const jobHandlers = `group:"job_handlers"`

func App() fx.Option {
	return fx.Options(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			l := &fxevent.ZapLogger{Logger: log.Named("fx")}
			l.UseLogLevel(zapcore.DebugLevel)
			return l
		}),
		fx.Provide(
			config.Load,
			func(c config.Config) config.Worker { return c.Worker },
			newLogger,
			newPostgres,
			newTxManager,
			queries.Load,
			fx.Annotate(clock.New,
				fx.As(fx.Self()),
				fx.As(new(intake.Clock)),
				fx.As(new(review.Clock)),
				fx.As(new(worker.Clock)),
				fx.As(new(httptransport.DemoClock)),
			),
		),
		fx.Provide(
			fx.Annotate(postgres.NewPatients, fx.As(new(intake.Patients)), fx.As(new(cases.Patients))),
			fx.Annotate(postgres.NewCases, fx.As(new(intake.Cases)), fx.As(new(cases.Cases)), fx.As(new(review.Cases))),
			fx.Annotate(postgres.NewRecommendations,
				fx.As(new(intake.Recommendations)),
				fx.As(new(cases.Recommendations)),
				fx.As(new(review.Recommendations)),
			),
			fx.Annotate(postgres.NewEvents, fx.As(new(intake.Events)), fx.As(new(review.Events))),
			fx.Annotate(postgres.NewJobs, fx.As(new(intake.Jobs)), fx.As(new(worker.Queue))),
			newMIS,
			newAIService,
		),
		fx.Provide(
			fx.Annotate(intake.NewService, fx.As(new(httptransport.Intake))),
			fx.Annotate(intake.NewAssessor, fx.As(new(worker.Handler)), fx.ResultTags(jobHandlers)),
			fx.Annotate(cases.NewService, fx.As(new(httptransport.Cases))),
			fx.Annotate(review.NewService, fx.As(new(httptransport.Review))),
			httptransport.NewHandler,
			fx.Annotate(worker.NewPool, fx.ParamTags(``, ``, ``, ``, jobHandlers)),
		),
		fx.Invoke(runHTTPServer, runWorkerPool),
	)
}

func newLogger(cfg config.Config) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	zcfg := zap.NewProductionConfig()
	zcfg.Level = zap.NewAtomicLevelAt(level)
	return zcfg.Build()
}

func newPostgres(lc fx.Lifecycle, cfg config.Config, log *zap.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Postgres.ConnectTimeout)
	defer cancel()
	pool, err := postgres.Connect(ctx, cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			applied, err := migrations.Up(ctx, pool)
			if err != nil {
				return err
			}
			log.Info("database migrated", zap.Strings("applied", applied))
			return nil
		},
		OnStop: func(context.Context) error {
			pool.Close()
			return nil
		},
	})
	return pool, nil
}

func newTxManager(pool *pgxpool.Pool) (trm.Manager, error) {
	return manager.New(trmpgx.NewDefaultFactory(pool))
}

func newAIService(cfg config.Config, log *zap.Logger) intake.AIService {
	if cfg.AIService.Mock {
		log.Warn("ai-service mock enabled")
		return aiservice.NewMock()
	}
	return aiservice.NewClient(cfg.AIService.URL, cfg.AIService.Timeout)
}

func newMIS(cfg config.Config) intake.MIS {
	return mis.NewClient(cfg.MIS.URL, cfg.MIS.Timeout)
}

func runHTTPServer(lc fx.Lifecycle, cfg config.Config, h *httptransport.Handler, log *zap.Logger) {
	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           h.Routes(cfg.DemoMode),
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("http server stopped", zap.Error(err))
				}
			}()
			log.Info("http server started", zap.String("addr", ln.Addr().String()), zap.Bool("demo_mode", cfg.DemoMode))
			return nil
		},
		OnStop: srv.Shutdown,
	})
}

func runWorkerPool(lc fx.Lifecycle, pool *worker.Pool) {
	lc.Append(fx.StartStopHook(pool.Start, pool.Stop))
}
