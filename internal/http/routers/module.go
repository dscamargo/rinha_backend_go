package routers

import "go.uber.org/fx"

var Module = fx.Provide(
	NewPessoaRouter,
	MakeRouter,
)
