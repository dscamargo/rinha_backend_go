package pessoasdb

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/dscamargo/rinha_backend_go/internal/domain/pessoa"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/rueidis"
	"os"
	"time"
)

type PessoasDbCache struct {
	client rueidis.Client
}

func NewPessoasDbCache() *PessoasDbCache {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{os.Getenv("CACHE_URL")},
		SelectDB:    0,
	})

	if err != nil {
		panic(err)
	}

	return &PessoasDbCache{client: client}
}

func (c *PessoasDbCache) GetPessoaCache(id string) (*pessoa.Pessoa, error) {
	cmd := c.client.B().Get().Key(id).Cache()
	pBytes, err := c.client.DoCache(context.Background(), cmd, 30*time.Minute).AsBytes()
	if err != nil {
		return nil, err
	}

	var ps pessoa.Pessoa
	err = sonic.Unmarshal(pBytes, &ps)
	if err != nil {
		return nil, err
	}

	return &ps, nil
}

func (c *PessoasDbCache) GetApelidoUtilizado(apelido string) (bool, error) {
	cmd := c.client.B().Getbit().Key(apelido).Offset(0).Cache()
	exists, err := c.client.DoCache(context.Background(), cmd, 30*time.Minute).AsBool()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (c *PessoasDbCache) SetPessoaEApelido(value *pessoa.Pessoa) error {
	valueStr, err := sonic.MarshalString(value)
	if err != nil {
		log.Errorf("Erro ao serializar pessoa %v", err)
		return err
	}

	cmdPessoa := c.client.B().Set().Key(value.ID).Value(valueStr).Build()
	cmdApelido := c.client.B().Setbit().Key(value.Apelido).Offset(0).Value(1).Build()

	cmds := make(rueidis.Commands, 0, 2)
	cmds = append(cmds, cmdPessoa)
	cmds = append(cmds, cmdApelido)

	for _, res := range c.client.DoMulti(context.Background(), cmds...) {
		err := res.Error()
		if err != nil {
			log.Errorf("Erro ao enviar DoMulti: %v", err)
			return err
		}
	}

	return nil
}
