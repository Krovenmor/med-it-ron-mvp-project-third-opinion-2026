package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/demo"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/http"
)

var demoModule = fx.Module("demo",
	fx.Provide(
		func(r *postgres.Demo) demo.Storage { return r },
		func(c *clock.Clock) demo.Clock { return c },
		fx.Annotate(demo.NewService, fx.As(new(httptransport.Demo))),
	),
)
