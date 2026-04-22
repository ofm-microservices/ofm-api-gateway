package appfx

import "go.uber.org/fx"

// Module aggregates the full FX graph for api-gateway.
var Module = fx.Options(
	ConfigModule,
	LoggerModule,
	AppModule,
	MessagingModule,
	ServiceModule,
	HTTPModule,
	HTTPV1Module,
)
