package pessoa

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/rueidis"
	"strings"
)

type Service struct {
	repository Repository
	cache      CacheRepository
}

func NewService(repository Repository, cache CacheRepository) *Service {
	return &Service{repository, cache}
}

func (s *Service) Create(nome, apelido, nascimento string, stack []string) (*Pessoa, error) {
	apelidoUtilizado, err := s.cache.GetApelidoUtilizado(apelido)
	if err != nil {
		return nil, err
	}

	if apelidoUtilizado {
		return nil, ErrApelidoJaUtilizado
	}

	pessoa := NewPessoa(nome, apelido, nascimento, stack)

	err = s.repository.Create(pessoa)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SetPessoaEApelido(pessoa); err != nil {
		log.Errorf("Erro ao salvar no redis", err)
		return nil, err
	}

	return pessoa, nil
}

func (s *Service) Count() (int64, error) {
	total, err := s.repository.Count()
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Service) FindById(id string) (*Pessoa, error) {
	cachedPessoa, err := s.cache.GetPessoaCache(id)

	if err != nil && !rueidis.IsRedisNil(err) {
		log.Errorf("Erro ao ler apelido do redis: %v", err)
		return nil, err
	}

	if cachedPessoa != nil {
		return cachedPessoa, nil
	}

	p, err := s.repository.GetById(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Search(term string) ([]Pessoa, error) {
	normalizedTerm := strings.ToLower(term)

	cachedResult, err := s.cache.GetSearch(normalizedTerm)
	if err != nil && !rueidis.IsRedisNil(err) {
		log.Errorf("Erro ao buscar pessoas no redis %v", err)
		return nil, err
	}

	if len(cachedResult) > 0 {
		return cachedResult, nil
	}

	result, err := s.repository.Search(term)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		go func() {
			if err := s.cache.SetSearch(normalizedTerm, result); err != nil {
				log.Errorf("Erro ao salvar no cache", err)
			}
		}()
	}

	return result, nil
}
