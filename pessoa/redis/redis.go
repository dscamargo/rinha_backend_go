package redis

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/dscamargo/rinha_backend_go/config"
	"github.com/dscamargo/rinha_backend_go/pessoa"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/rueidis"
	"time"
)

type CacheRepository struct {
	client rueidis.Client
}

func NewCacheRepository(cfg *config.Config) *CacheRepository {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{cfg.Cache.Url},
		SelectDB:    0,
	})

	if err != nil {
		panic(err)
	}

	return &CacheRepository{client: client}
}

func (r *CacheRepository) GetPessoaCache(id string) (*pessoa.Pessoa, error) {
	cmd := r.client.B().Get().Key(id).Cache()
	pBytes, err := r.client.DoCache(context.Background(), cmd, 30*time.Minute).AsBytes()

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

func (r *CacheRepository) GetApelidoUtilizado(apelido string) (bool, error) {
	cmd := r.client.B().Getbit().Key(apelido).Offset(0).Cache()
	exists, err := r.client.DoCache(context.Background(), cmd, 30*time.Minute).AsBool()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (r *CacheRepository) SetPessoaEApelido(value *pessoa.Pessoa) error {
	valueStr, err := sonic.MarshalString(value)
	if err != nil {
		log.Errorf("Erro ao serializar pessoa %v", err)
		return err
	}

	cmdPessoa := r.client.B().Set().Key(value.ID).Value(valueStr).Ex(1 * time.Minute).Build()
	cmdApelido := r.client.B().Setbit().Key(value.Apelido).Offset(0).Value(1).Build()

	cmds := make(rueidis.Commands, 0, 2)
	cmds = append(cmds, cmdPessoa)
	cmds = append(cmds, cmdApelido)

	for _, res := range r.client.DoMulti(context.Background(), cmds...) {
		err := res.Error()
		if err != nil {
			log.Errorf("Erro ao enviar DoMulti: %v", err)
			return err
		}
	}

	return nil
}

func (r *CacheRepository) SetSearch(term string, result []pessoa.Pessoa) error {
	valStr, err := sonic.MarshalString(result)
	if err != nil {
		return err
	}
	cmd := r.client.B().Set().Key("busca:" + term).Value(valStr).Ex(15 * time.Minute).Build()
	return r.client.Do(context.Background(), cmd).Error()
}

func (r *CacheRepository) GetSearch(term string) ([]pessoa.Pessoa, error) {
	cmd := r.client.B().Get().Key("busca:" + term).Cache()
	pBytes, err := r.client.DoCache(context.Background(), cmd, 15*time.Minute).AsBytes()
	if err != nil {
		return nil, err
	}

	var ps []pessoa.Pessoa
	err = sonic.Unmarshal(pBytes, &ps)
	if err != nil {
		return nil, err
	}

	return ps, nil
}
