package pessoa

import "go.uber.org/fx"

var Module = fx.Provide(
	NewCountPessoas,
	NewFindPessoa,
	NewCreatePessoa,
)
