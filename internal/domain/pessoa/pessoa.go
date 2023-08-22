package pessoa

import (
	"github.com/google/uuid"
	"strings"
)

type Pessoa struct {
	ID         string
	Nome       string
	Apelido    string
	Nascimento string
	Stack      []string
	Search     string
}

func NewPessoa(nome, apelido, nascimento string, stack []string) *Pessoa {
	return &Pessoa{
		ID:         uuid.NewString(),
		Nome:       nome,
		Apelido:    apelido,
		Nascimento: nascimento,
		Stack:      stack,
	}
}

func (p *Pessoa) StackStr() string {
	return strings.Join(p.Stack, ",")
}

func (p *Pessoa) SearchStr() string {
	return p.Nome + " " + p.Apelido + " " + p.StackStr()
}
