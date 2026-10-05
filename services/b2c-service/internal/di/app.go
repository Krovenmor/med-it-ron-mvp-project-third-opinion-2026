package di

import "go.uber.org/fx"

func App() fx.Option {
	return fx.Options(
		fx.WithLogger(newFxLogger),
		coreModule,
		clientsModule,
		routeModule,
		bookingModule,
		stepsModule,
		httpModule,
	)
}
