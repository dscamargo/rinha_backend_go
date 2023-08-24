package redis

import (
	"github.com/dscamargo/rinha_backend_go/pessoa"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(NewCacheRepository,
		fx.As(new(pessoa.CacheRepository)),
	),
)
