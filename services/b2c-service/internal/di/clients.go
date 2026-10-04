package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/infra/b2b"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/infra/mis"
)

var clientsModule = fx.Module("clients",
	fx.Provide(
		newB2BClient,
		newMISClient,
	),
)

func newB2BClient(cfg config.Config) *b2b.Client {
	return b2b.NewClient(cfg.B2B.URL, cfg.B2B.Timeout, cfg.MIS.SourceSystem)
}

func newMISClient(cfg config.Config) *mis.Client {
	return mis.NewClient(cfg.MIS.URL, cfg.MIS.Timeout)
}
