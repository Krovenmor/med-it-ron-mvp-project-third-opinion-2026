package di

import "go.uber.org/fx"

func App() fx.Option {
	return fx.Options(
		fx.WithLogger(newFxLogger),
		coreModule,
		postgresModule,
		intakeModule,
		casesModule,
		reviewModule,
		planModule,
		bookingModule,
		httpModule,
		workerModule,
	)
}
