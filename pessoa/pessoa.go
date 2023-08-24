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

type Repository interface {
	Create(p *Pessoa) error
	GetById(id string) (*Pessoa, error)
	Search(term string) ([]Pessoa, error)
	Count() (int64, error)
}

type CacheRepository interface {
	GetPessoaCache(id string) (*Pessoa, error)
	GetApelidoUtilizado(apelido string) (bool, error)
	SetPessoaEApelido(value *Pessoa) error
	SetSearch(term string, result []Pessoa) error
	GetSearch(term string) ([]Pessoa, error)
}

type UseCase interface {
	Create(nome, apelido, nascimento string, stack []string) (*Pessoa, error)
	Count() (int64, error)
	FindById(id string) (*Pessoa, error)
	Search(term string) ([]Pessoa, error)
}
