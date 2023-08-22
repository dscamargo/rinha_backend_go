package pessoasdb

import (
	"github.com/dscamargo/rinha_backend_go/internal/domain/pessoa"
	"go.uber.org/fx"
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
