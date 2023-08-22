package pessoa

type CountPessoas struct {
	repository Repository
}

func NewCountPessoas(repository Repository) *CountPessoas {
	return &CountPessoas{repository}
}

func (c *CountPessoas) Exec() (int64, error) {
	total, err := c.repository.Count()
	if err != nil {
		return 0, err
	}
	return total, nil
}
