package di

import (
	"context"
	"errors"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/config"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/transport/http"
)

var httpModule = fx.Module("http",
	fx.Provide(httptransport.NewHandler),
	fx.Invoke(runHTTPServer),
)

func runHTTPServer(lc fx.Lifecycle, cfg config.Config, h *httptransport.Handler, log *zap.Logger) {
	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           h.Routes(),
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
			log.Info("http server started", zap.String("addr", ln.Addr().String()))
			return nil
		},
		OnStop: srv.Shutdown,
	})
}
