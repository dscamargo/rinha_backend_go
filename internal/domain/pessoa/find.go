package pessoa

type FindPessoa struct {
	repository Repository
}

func NewFindPessoa(repository Repository) *FindPessoa {
	return &FindPessoa{repository}
}

func (f *FindPessoa) Exec(id string) (*Pessoa, error) {
	p, err := f.repository.GetById(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}
