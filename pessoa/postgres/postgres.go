package postgres

import (
	"context"
	"errors"
	"github.com/dscamargo/rinha_backend_go/pessoa"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type PGRepository struct {
	db       *pgxpool.Pool
	jobQueue JobQueue
}

func NewPGRepository(db *pgxpool.Pool, jobQueue JobQueue) *PGRepository {
	return &PGRepository{db, jobQueue}
}

func (r *PGRepository) GetById(id string) (*pessoa.Pessoa, error) {
	var p pessoa.Pessoa
	var stackStr string

	err := r.db.QueryRow(context.Background(), QuerySelectPessoaById, id).Scan(&p.ID, &p.Nome, &p.Apelido, &p.Nascimento, &stackStr)

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

func (r *PGRepository) Search(term string) ([]pessoa.Pessoa, error) {
	rows, err := r.db.Query(context.Background(), QuerySelectPessoasByTerm, term)
	if err != nil {
		log.Errorf("Erro ao buscar pessoas no banco de dados %v", err)
		return nil, err
	}
	defer rows.Close()

	result, err := formatPessoasResult(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PGRepository) Count() (int64, error) {
	var total int64
	err := r.db.QueryRow(context.Background(), QueryCountAllPessoas).Scan(&total)
	if err != nil {
		log.Errorf("Erro ao buscar pessoas no banco de dados %v", err)
		return 0, err
	}
	return total, nil
}

func (r *PGRepository) Create(p *pessoa.Pessoa) error {
	r.jobQueue <- Job{Payload: p}
	return nil
}

func formatPessoasResult(rows pgx.Rows) ([]pessoa.Pessoa, error) {
	result := make([]pessoa.Pessoa, 0)
	for rows.Next() {
		var ps pessoa.Pessoa
		var stackStr string
		err := rows.Scan(&ps.ID, &ps.Nome, &ps.Apelido, &ps.Nascimento, &stackStr)
		if err != nil {
			log.Errorf("Erro ao buscar pessoas no banco de dados %v", err)
			return nil, err
		}

		ps.Stack = strings.Split(stackStr, ",")

		result = append(result, ps)
	}
	return result, nil
}
