package db

import (
	"github.com/dscamargo/rinha_backend_go/internal/infra/db/pessoasdb"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewPGDatabase,
	),
	pessoasdb.Module,
)
