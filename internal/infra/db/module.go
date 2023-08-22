package db

import (
	"go.uber.org/fx"
	"rinha_v2/internal/infra/db/pessoasdb"
)

var Module = fx.Options(
	fx.Provide(
		NewPGDatabase,
	),
	pessoasdb.Module,
)
