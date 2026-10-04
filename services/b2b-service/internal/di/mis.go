package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/mis"
)

var misModule = fx.Module("mis",
	fx.Provide(newMISClient),
)

func newMISClient(cfg config.Config) *mis.Client {
	return mis.NewClient(cfg.MIS.URL, cfg.MIS.Timeout)
}
