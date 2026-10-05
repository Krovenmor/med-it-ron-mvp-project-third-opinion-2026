package system

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	careapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/api"
	doctorapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/api"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/service/demo"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/transport/http"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
)

var Module = fx.Module("system",
	fx.Provide(
		fx.Private,
		newDemoService,
	),
	fx.Provide(
		fx.Annotate(newHandler, fx.As(new(httpx.Routes)), fx.ResultTags(`group:"http_routes"`)),
	),
)

func newDemoService(c *clock.Clock, gateway gatewayapi.Demo, doctor doctorapi.Demo, care careapi.Demo) *demo.Service {
	return demo.NewService(c, gateway, doctor, care)
}

func newHandler(c *clock.Clock, d *demo.Service, cfg config.Config, log *zap.Logger) *httptransport.Handler {
	return httptransport.NewHandler(c, d, cfg.DemoMode, log)
}
