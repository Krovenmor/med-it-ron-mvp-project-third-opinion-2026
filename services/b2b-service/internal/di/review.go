package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/review"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/http"
)

var reviewModule = fx.Module("review",
	fx.Provide(
		func(c *clock.Clock) review.Clock { return c },
		func(r *postgres.Cases) review.Cases { return r },
		func(r *postgres.Recommendations) review.Recommendations { return r },
		func(r *postgres.Events) review.Events { return r },
		fx.Annotate(review.NewService, fx.As(new(httptransport.Review))),
	),
)
