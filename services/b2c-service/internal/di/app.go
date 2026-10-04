package di

import "go.uber.org/fx"

func App() fx.Option {
	return fx.Options(
		fx.WithLogger(newFxLogger),
		coreModule,
		clientsModule,
		graphModule,
		bookingModule,
		httpModule,
	)
}
