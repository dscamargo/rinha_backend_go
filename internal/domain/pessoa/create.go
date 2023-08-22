package pessoa

type CreatePessoa struct {
	repository Repository
}

func NewCreatePessoa(repository Repository) *CreatePessoa {
	return &CreatePessoa{repository}
}

func (c *CreatePessoa) Exec(nome, apelido, nascimento string, stack []string) (*Pessoa, error) {
	apelidoUtilizado, err := c.repository.CheckApelido(apelido)
	if err != nil {
		return nil, err
	}

	if apelidoUtilizado {
		return nil, ErrApelidoJaUtilizado
	}

	pessoa := NewPessoa(nome, apelido, nascimento, stack)

	err = c.repository.Create(pessoa)
	if err != nil {
		return nil, err
	}

	return pessoa, nil
}
