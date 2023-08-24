package pessoa

import "errors"

var (
	ErrPessoaNotFound     = errors.New("pessoa not found")
	ErrApelidoJaUtilizado = errors.New("apelido já utilizado")
	ErrInvalidBody        = errors.New("invalid body")
)
