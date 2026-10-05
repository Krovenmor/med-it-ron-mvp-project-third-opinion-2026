package di

import (
	"context"
	"errors"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/events"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
)

var platformModule = fx.Module("platform",
	fx.Provide(
		fx.Annotate(events.NewBus, fx.ParamTags(`group:"event_subscriptions"`)),
	),
	fx.Invoke(fx.Annotate(runHTTPServer, fx.ParamTags(``, ``, `group:"http_routes"`, ``))),
)

func runHTTPServer(lc fx.Lifecycle, cfg config.Config, routes []httpx.Routes, log *zap.Logger) {
	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           httpx.NewMux(routes, httpx.NewResponder(log, "server")),
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("http server stopped", zap.Error(err))
				}
			}()
			log.Info("http server started", zap.String("addr", ln.Addr().String()), zap.Bool("demo_mode", cfg.DemoMode))
			return nil
		},
		OnStop: srv.Shutdown,
	})
}
