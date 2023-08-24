package pessoa

import "go.uber.org/fx"

var Module = fx.Provide(
	fx.Annotate(NewService,
		fx.As(new(UseCase)),
	),
)
