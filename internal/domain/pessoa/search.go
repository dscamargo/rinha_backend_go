package pessoa

type SearchPessoas struct {
	repository Repository
}

func NewSearchPessoas(repository Repository) *SearchPessoas {
	return &SearchPessoas{repository}
}

func (f *SearchPessoas) Exec(term string) ([]Pessoa, error) {
	p, err := f.repository.Search(term)
	if err != nil {
		return nil, err
	}
	return p, nil
}
