package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/infra/b2b"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/service/steps"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/transport/http"
)

var stepsModule = fx.Module("steps",
	fx.Provide(
		func(c *b2b.Client) steps.Plans { return c },
		fx.Annotate(steps.NewService, fx.As(new(httptransport.Steps))),
	),
)
