package controllers

import "github.com/dscamargo/rinha_backend_go/internal/domain/pessoa"

type PessoaOutput struct {
	ID         string   `json:"id"`
	Apelido    string   `json:"apelido"`
	Nome       string   `json:"nome"`
	Nascimento string   `json:"nascimento"`
	Stack      []string `json:"stack"`
}

func mapPessoaOutput(p *pessoa.Pessoa) PessoaOutput {
	return PessoaOutput{
		ID:         p.ID,
		Apelido:    p.Apelido,
		Nome:       p.Nome,
		Nascimento: p.Nascimento,
		Stack:      p.Stack,
	}
}
