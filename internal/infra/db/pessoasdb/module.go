package pessoasdb

import (
	"go.uber.org/fx"
	"rinha_v2/internal/domain/pessoa"
)

var Module = fx.Provide(
	NewPessoasDbCache,
	NewWorker,
	NewInsertChannel,
	fx.Annotate(
		NewPessoaRepository,
		fx.As(new(pessoa.Repository)),
	),
)
