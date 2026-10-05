package care

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/booking"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/dashboard"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/demo"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/operator"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/patient"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/plan"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/protocol"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/routes"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/transport/http"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/events"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
	platformpg "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

const name = "care"

var Module = fx.Module(name,
	fx.Provide(
		fx.Private,
		queries.Load,
		openDatabase,
		func(db *platformpg.Database) trm.Manager { return db.Tx },
		jobs.NewQueue,
		events.NewInbox,
		newConsumer,
		postgres.NewPatients,
		postgres.NewRoutes,
		postgres.NewPlanItems,
		postgres.NewBookings,
		postgres.NewOperatorTasks,
		postgres.NewCallAttempts,
		postgres.NewNotifications,
		postgres.NewEvents,
		postgres.NewDashboard,
		postgres.NewDemo,
	),
	fx.Provide(
		fx.Private,
		func(c *clock.Clock) routes.Clock { return c },
		func(r *postgres.Patients) routes.Patients { return r },
		func(r *postgres.Routes) routes.Routes { return r },
		func(r *postgres.PlanItems) routes.PlanItems { return r },
		func(r *postgres.Notifications) routes.Notifications { return r },
		func(r *postgres.Events) routes.Events { return r },
		func(p *protocol.Planner) routes.Protocol { return p },
		routes.NewService,
	),
	fx.Provide(
		fx.Private,
		func(c *clock.Clock) protocol.Clock { return c },
		func(r *postgres.Routes) protocol.Routes { return r },
		func(r *postgres.Patients) protocol.Patients { return r },
		func(r *postgres.OperatorTasks) protocol.Tasks { return r },
		func(q *jobs.Queue) protocol.Jobs { return q },
		func(r *postgres.Notifications) protocol.Notifications { return r },
		func(r *postgres.Events) protocol.Events { return r },
		protocol.NewPlanner,
		protocol.NewHandlers,
	),
	fx.Provide(
		fx.Private,
		func(c *clock.Clock) booking.Clock { return c },
		func(r *postgres.Routes) booking.Routes { return r },
		func(r *postgres.PlanItems) booking.PlanItems { return r },
		func(r *postgres.Bookings) booking.Bookings { return r },
		func(r *postgres.Events) booking.Events { return r },
		func(p *protocol.Planner) booking.Protocol { return p },
		booking.NewService,
	),
	fx.Provide(
		fx.Private,
		func(c *clock.Clock) operator.Clock { return c },
		func(r *postgres.OperatorTasks) operator.Tasks { return r },
		func(r *postgres.CallAttempts) operator.Attempts { return r },
		func(r *postgres.Routes) operator.Routes { return r },
		func(r *postgres.Patients) operator.Patients { return r },
		func(r *postgres.PlanItems) operator.PlanItems { return r },
		func(r *postgres.Bookings) operator.Bookings { return r },
		func(r *postgres.Notifications) operator.Notifications { return r },
		func(r *postgres.Events) operator.Events { return r },
		func(p *protocol.Planner) operator.Protocol { return p },
		func(m gatewayapi.MIS) operator.MIS { return m },
		func(s *booking.Service) operator.Booker { return s },
		operator.NewService,
	),
	fx.Provide(
		fx.Private,
		func(c *clock.Clock) plan.Clock { return c },
		func(r *postgres.Patients) plan.Patients { return r },
		func(r *postgres.Routes) plan.Routes { return r },
		func(r *postgres.PlanItems) plan.PlanItems { return r },
		plan.NewService,
		func(c *clock.Clock) dashboard.Clock { return c },
		func(r *postgres.Dashboard) dashboard.Facts { return r },
		dashboard.NewService,
		func(r *postgres.Demo) demo.Storage { return r },
		func(c *clock.Clock) patient.Clock { return c },
		func(r *postgres.Routes) patient.Routes { return r },
		func(r *postgres.PlanItems) patient.PlanItems { return r },
		func(r *postgres.Bookings) patient.Bookings { return r },
		func(r *postgres.OperatorTasks) patient.Tasks { return r },
		func(r *postgres.Events) patient.Events { return r },
		func(p *protocol.Planner) patient.Protocol { return p },
		patient.NewService,
		func(s *patient.Service) httptransport.Patient { return s },
		func(s *plan.Service) httptransport.Plan { return s },
		func(s *booking.Service) httptransport.Bookings { return s },
		func(s *operator.Service) httptransport.Operator { return s },
		func(s *dashboard.Service) httptransport.Dashboard { return s },
		newWorkerPool,
	),
	fx.Provide(
		func(s *routes.Service) api.Routes { return s },
		func(tx trm.Manager, storage demo.Storage) api.Demo { return demo.NewService(tx, storage) },
		fx.Annotate(httptransport.NewHandler, fx.As(new(httpx.Routes)), fx.ResultTags(`group:"http_routes"`)),
		fx.Annotate(
			func(s *routes.Service, c *events.Consumer) []events.Subscription { return s.Subscriptions(c) },
			fx.ResultTags(`group:"event_subscriptions,flatten"`),
		),
	),
	fx.Invoke(jobs.Run),
)

func openDatabase(lc fx.Lifecycle, log *zap.Logger, cfg config.Config, q queries.Queries) (*platformpg.Database, error) {
	return platformpg.OpenModule(lc, log, cfg.Postgres, postgres.Schema(q))
}

func newConsumer(tx trm.Manager, inbox *events.Inbox, c *clock.Clock) *events.Consumer {
	return events.NewConsumer(tx, inbox, c)
}

func newWorkerPool(cfg config.Config, queue *jobs.Queue, c *clock.Clock, log *zap.Logger, h *protocol.Handlers) (*jobs.Pool, error) {
	return jobs.NewPool(name, cfg.Worker, queue, c, log, h.All())
}
