package pessoasdb

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/rueidis"
	"rinha_v2/internal/domain/pessoa"
	"strings"
)

type PessoaRepository struct {
	db         *pgxpool.Pool
	cache      *PessoasDbCache
	insertChan chan pessoa.Pessoa
}

func NewPessoaRepository(db *pgxpool.Pool, cache *PessoasDbCache, insertChannel chan pessoa.Pessoa) pessoa.Repository {
	return &PessoaRepository{db, cache, insertChannel}
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
	rows, err := r.db.Query(context.Background(), QuerySelectPessoasByTerm, strings.ToLower(term))
	if err != nil {
		log.Errorf("Erro ao buscar pessoas no banco de dados %v", err)
		return nil, err
	}
	return formatPessoasResult(rows)
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

func (r *PessoaRepository) Create(p *pessoa.Pessoa) (*pessoa.Pessoa, error) {
	apelidoJaUtilizado, err := r.cache.GetApelidoUtilizado(p.Apelido)
	if err != nil {
		log.Errorf("Erro ao buscar apelido no redis %v", err)
		return nil, err
	}

	if apelidoJaUtilizado {
		return nil, pessoa.ErrApelidoJaUtilizado
	}

	r.insertChan <- *p

	err = r.cache.SetPessoaEApelido(p)
	if err != nil {
		log.Errorf("Erro ao salvar no redis", err)
	}

	return p, nil
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
