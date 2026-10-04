package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/worker"
)

const jobHandlers = `group:"job_handlers"`

var workerModule = fx.Module("worker",
	fx.Provide(
		func(c config.Config) config.Worker { return c.Worker },
		func(c *clock.Clock) worker.Clock { return c },
		func(r *postgres.Jobs) worker.Queue { return r },
		fx.Annotate(worker.NewPool, fx.ParamTags(``, ``, ``, ``, jobHandlers)),
	),
	fx.Invoke(runWorkerPool),
)

func runWorkerPool(lc fx.Lifecycle, pool *worker.Pool) {
	lc.Append(fx.StartStopHook(pool.Start, pool.Stop))
}
