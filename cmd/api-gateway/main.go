package main

import (
	appfx "api-gateway/internal/fx"

	"go.uber.org/fx"
)

type runner interface {
	Run()
}

var newApp = func(opts ...fx.Option) runner {
	return fx.New(opts...)
}

func main() {
	app := newApp(
		appfx.ConfigModule,
		appfx.LoggerModule,
		appfx.TracingModule,
		appfx.MetricsModule,
		appfx.AppModule,
		appfx.MessagingModule,
		appfx.ServiceModule,
		appfx.HTTPModule,
		appfx.HTTPV1Module,
	)
	app.Run()
}
