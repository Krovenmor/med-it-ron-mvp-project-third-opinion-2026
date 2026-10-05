package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/infra/b2b"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/infra/mis"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/service/route"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/transport/http"
)

var routeModule = fx.Module("route",
	fx.Provide(
		func(c config.Config) config.Clinic { return c.Clinic },
		func(c *b2b.Client) route.Plans { return c },
		func(c *mis.Client) route.MIS { return c },
		fx.Annotate(route.NewService, fx.As(new(httptransport.Routes))),
	),
)
