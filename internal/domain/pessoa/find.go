package pessoa

type FindPessoa struct {
	repository Repository
}

func NewFindPessoa(repository Repository) *FindPessoa {
	return &FindPessoa{repository}
}

func (f *FindPessoa) FindById(id string) (*Pessoa, error) {
	p, err := f.repository.GetById(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (f *FindPessoa) Search(term string) ([]Pessoa, error) {
	p, err := f.repository.Search(term)
	if err != nil {
		return nil, err
	}
	return p, nil
}
