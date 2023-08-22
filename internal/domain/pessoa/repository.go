package pessoa

type Repository interface {
	Create(p *Pessoa) (*Pessoa, error)
	GetById(id string) (*Pessoa, error)
	Search(term string) ([]Pessoa, error)
	Count() (int64, error)
}
