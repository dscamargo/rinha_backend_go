package pessoasdb

import (
	"context"
	"errors"
	"github.com/dscamargo/rinha_backend_go/internal/domain/pessoa"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/rueidis"
	"strings"
)

type PessoaRepository struct {
	db       *pgxpool.Pool
	cache    *PessoasDbCache
	jobQueue JobQueue
}

func NewPessoaRepository(db *pgxpool.Pool, cache *PessoasDbCache, jobQueue JobQueue) pessoa.Repository {
	return &PessoaRepository{db, cache, jobQueue}
}

func (r *PessoaRepository) GetById(id string) (*pessoa.Pessoa, error) {
	cachedPessoa, err := r.cache.GetPessoaCache(id)

	if err != nil && !rueidis.IsRedisNil(err) {
		log.Errorf("Erro ao ler apelido do redis: %v", err)
		return nil, err
	}

	if cachedPessoa != nil {
		return cachedPessoa, nil
	}

	var p pessoa.Pessoa
	var stackStr string

	err = r.db.QueryRow(context.Background(), QuerySelectPessoaById, id).Scan(&p.ID, &p.Nome, &p.Apelido, &p.Nascimento, &stackStr)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pessoa.ErrPessoaNotFound
		}
		log.Errorf("Erro ao buscar pessoa no banco de dados %v", err)
		return nil, err
	}

	if len(stackStr) > 0 {
		p.Stack = strings.Split(stackStr, ",")
	}

	return &p, nil
}

func (r *PessoaRepository) Search(term string) ([]pessoa.Pessoa, error) {
	normalizedTerm := strings.ToLower(term)
	cachedResult, err := r.cache.GetSearch(normalizedTerm)
	if err != nil && !rueidis.IsRedisNil(err) {
		log.Errorf("Erro ao buscar pessoas no redis %v", err)
		return nil, err
	}

	if len(cachedResult) > 0 {
		return cachedResult, nil
	}

	rows, err := r.db.Query(context.Background(), QuerySelectPessoasByTerm, normalizedTerm)
	if err != nil {
		log.Errorf("Erro ao buscar pessoas no banco de dados %v", err)
		return nil, err
	}
	defer rows.Close()

	result, err := formatPessoasResult(rows)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		go func() {
			if err := r.cache.SetSearch(normalizedTerm, result); err != nil {
				log.Errorf("Erro ao salvar no cache", err)
			}
		}()
	}
	return result, nil
}

func (r *PessoaRepository) Count() (int64, error) {
	var total int64
	err := r.db.QueryRow(context.Background(), QueryCountAllPessoas).Scan(&total)
	if err != nil {
		log.Errorf("Erro ao buscar pessoas no banco de dados %v", err)
		return 0, err
	}
	return total, nil
}

func (r *PessoaRepository) Create(p *pessoa.Pessoa) error {
	if err := r.cache.SetPessoaEApelido(p); err != nil {
		log.Errorf("Erro ao salvar no redis", err)
		return err
	}

	r.jobQueue <- Job{Payload: p}

	return nil
}

func (r *PessoaRepository) CheckApelido(apelido string) (bool, error) {
	apelidoJaUtilizado, err := r.cache.GetApelidoUtilizado(apelido)
	if err != nil {
		log.Errorf("Erro ao buscar apelido no redis %v", err)
		return false, err
	}
	return apelidoJaUtilizado, nil
}

func formatPessoasResult(rows pgx.Rows) ([]pessoa.Pessoa, error) {
	result := make([]pessoa.Pessoa, 0)
	for rows.Next() {
		var person pessoa.Pessoa
		var stackStr string
		err := rows.Scan(&person.ID, &person.Nome, &person.Apelido, &person.Nascimento, &stackStr)
		if err != nil {
			log.Errorf("Erro ao buscar pessoas no banco de dados %v", err)
			return nil, err
		}

		person.Stack = strings.Split(stackStr, ",")

		result = append(result, person)
	}
	return result, nil
}
