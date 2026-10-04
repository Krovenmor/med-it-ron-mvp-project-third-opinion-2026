package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/plan"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/http"
)

var planModule = fx.Module("plan",
	fx.Provide(
		func(r *postgres.Patients) plan.Patients { return r },
		func(r *postgres.Cases) plan.Cases { return r },
		func(r *postgres.Recommendations) plan.Recommendations { return r },
		fx.Annotate(plan.NewService, fx.As(new(httptransport.Plan))),
	),
)
