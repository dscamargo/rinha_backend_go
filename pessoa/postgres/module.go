package postgres

import (
	"github.com/dscamargo/rinha_backend_go/pessoa"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewDispatcher,
	NewJobQueue,
	fx.Annotate(NewPGRepository,
		fx.As(new(pessoa.Repository)),
	),
)
