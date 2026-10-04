package di

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/aiservice"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/mis"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/intake"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/http"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/worker"
)

var intakeModule = fx.Module("intake",
	fx.Provide(
		func(c *clock.Clock) intake.Clock { return c },
		func(r *postgres.Patients) intake.Patients { return r },
		func(r *postgres.Cases) intake.Cases { return r },
		func(r *postgres.Recommendations) intake.Recommendations { return r },
		func(r *postgres.Events) intake.Events { return r },
		func(r *postgres.Jobs) intake.Jobs { return r },
		newAIService,
		newMIS,
		fx.Annotate(intake.NewService, fx.As(new(httptransport.Intake))),
		fx.Annotate(intake.NewAssessor, fx.As(new(worker.Handler)), fx.ResultTags(jobHandlers)),
	),
)

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
