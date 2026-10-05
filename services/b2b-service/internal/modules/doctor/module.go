package doctor

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/service/cases"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/service/demo"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/service/escalation"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/service/intake"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/service/review"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/transport/http"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/events"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
	platformpg "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

const name = "doctor"

var Module = fx.Module(name,
	fx.Provide(
		fx.Private,
		queries.Load,
		openDatabase,
		func(db *platformpg.Database) trm.Manager { return db.Tx },
		jobs.NewQueue,
		events.NewInbox,
		newPublisher,
		newConsumer,
		postgres.NewCases,
		postgres.NewRecommendations,
		postgres.NewEvents,
		postgres.NewDemo,
		func(c *clock.Clock) review.Clock { return c },
		func(c *clock.Clock) escalation.Clock { return c },
		func(r *postgres.Cases) intake.Cases { return r },
		func(r *postgres.Cases) review.Cases { return r },
		func(r *postgres.Cases) cases.Cases { return r },
		func(r *postgres.Cases) escalation.Cases { return r },
		func(r *postgres.Recommendations) intake.Recommendations { return r },
		func(r *postgres.Recommendations) review.Recommendations { return r },
		func(r *postgres.Recommendations) cases.Recommendations { return r },
		func(r *postgres.Events) intake.Events { return r },
		func(r *postgres.Events) review.Events { return r },
		func(r *postgres.Events) escalation.Events { return r },
		func(q *jobs.Queue) intake.Jobs { return q },
		func(q *jobs.Queue) review.Jobs { return q },
		func(p *events.Publisher) intake.Publisher { return p },
		func(p *events.Publisher) review.Publisher { return p },
		func(p *events.Publisher) escalation.Publisher { return p },
		func(m gatewayapi.MIS) cases.MIS { return m },
		func(r *postgres.Demo) demo.Storage { return r },
		intake.NewService,
		review.NewService,
		cases.NewService,
		escalation.NewHandler,
		func(s *cases.Service) httptransport.Cases { return s },
		func(s *review.Service) httptransport.Review { return s },
		newWorkerPool,
	),
	fx.Provide(
		func(tx trm.Manager, storage demo.Storage) api.Demo { return demo.NewService(tx, storage) },
		fx.Annotate(httptransport.NewHandler, fx.As(new(httpx.Routes)), fx.ResultTags(`group:"http_routes"`)),
		fx.Annotate(
			func(s *intake.Service, c *events.Consumer) events.Subscription { return s.Subscription(c) },
			fx.ResultTags(`group:"event_subscriptions"`),
		),
	),
	fx.Invoke(jobs.Run),
)

func openDatabase(lc fx.Lifecycle, log *zap.Logger, cfg config.Config, q queries.Queries) (*platformpg.Database, error) {
	return platformpg.OpenModule(lc, log, cfg.Postgres, postgres.Schema(q))
}

func newPublisher(q *jobs.Queue, c *clock.Clock) *events.Publisher {
	return events.NewPublisher(q, c)
}

func newConsumer(tx trm.Manager, inbox *events.Inbox, c *clock.Clock) *events.Consumer {
	return events.NewConsumer(tx, inbox, c)
}

func newWorkerPool(
	cfg config.Config,
	queue *jobs.Queue,
	c *clock.Clock,
	log *zap.Logger,
	bus *events.Bus,
	escalate *escalation.Handler,
) (*jobs.Pool, error) {
	return jobs.NewPool(name, cfg.Worker, queue, c, log, []jobs.Handler{escalate, events.Delivery(name, bus)})
}
