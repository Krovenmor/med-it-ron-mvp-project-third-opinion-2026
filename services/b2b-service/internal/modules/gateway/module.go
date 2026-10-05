package gateway

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	careapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/infra/aiservice"
	misclient "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/infra/mis"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/service/demo"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/service/intake"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/service/mis"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/transport/http"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/events"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
	platformpg "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

const name = "gateway"

var Module = fx.Module(name,
	fx.Provide(
		fx.Private,
		queries.Load,
		openDatabase,
		func(db *platformpg.Database) trm.Manager { return db.Tx },
		jobs.NewQueue,
		newPublisher,
		postgres.NewPatients,
		postgres.NewIntakes,
		postgres.NewDemo,
		newMISClient,
		newAIService,
		func(c *clock.Clock) intake.Clock { return c },
		func(r *postgres.Patients) intake.Patients { return r },
		func(r *postgres.Intakes) intake.Intakes { return r },
		func(q *jobs.Queue) intake.Jobs { return q },
		func(p *events.Publisher) intake.Publisher { return p },
		func(c *misclient.Client) intake.MIS { return c },
		func(r careapi.Routes) intake.Routes { return r },
		func(r *postgres.Patients) mis.Patients { return r },
		func(c *misclient.Client) mis.Client { return c },
		func(r *postgres.Demo) demo.Storage { return r },
		intake.NewService,
		intake.NewAssessor,
		func(s *intake.Service) httptransport.Intake { return s },
		newWorkerPool,
	),
	fx.Provide(
		func(tx trm.Manager, storage demo.Storage) api.Demo { return demo.NewService(tx, storage) },
		func(patients mis.Patients, client mis.Client) api.MIS { return mis.NewService(patients, client) },
		fx.Annotate(httptransport.NewHandler, fx.As(new(httpx.Routes)), fx.ResultTags(`group:"http_routes"`)),
	),
	fx.Invoke(jobs.Run),
)

func openDatabase(lc fx.Lifecycle, log *zap.Logger, cfg config.Config, q queries.Queries) (*platformpg.Database, error) {
	return platformpg.OpenModule(lc, log, cfg.Postgres, postgres.Schema(q))
}

func newPublisher(q *jobs.Queue, c *clock.Clock) *events.Publisher {
	return events.NewPublisher(q, c)
}

func newMISClient(cfg config.Config) *misclient.Client {
	return misclient.NewClient(cfg.MIS.URL, cfg.MIS.Timeout)
}

func newAIService(cfg config.Config, log *zap.Logger) intake.AIService {
	if cfg.AIService.Mock {
		log.Warn("ai-service mock enabled")
		return aiservice.NewMock()
	}
	return aiservice.NewClient(cfg.AIService.URL, cfg.AIService.Timeout)
}

func newWorkerPool(
	cfg config.Config,
	queue *jobs.Queue,
	c *clock.Clock,
	log *zap.Logger,
	bus *events.Bus,
	assessor *intake.Assessor,
) (*jobs.Pool, error) {
	return jobs.NewPool(name, cfg.Worker, queue, c, log, []jobs.Handler{assessor, events.Delivery(name, bus)})
}
