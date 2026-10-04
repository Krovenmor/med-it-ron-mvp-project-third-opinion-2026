package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/mis"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/cases"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/http"
)

var casesModule = fx.Module("cases",
	fx.Provide(
		func(r *postgres.Cases) cases.Cases { return r },
		func(r *postgres.Patients) cases.Patients { return r },
		func(r *postgres.Recommendations) cases.Recommendations { return r },
		func(c *mis.Client) cases.MIS { return c },
		fx.Annotate(cases.NewService, fx.As(new(httptransport.Cases))),
	),
)
