package pessoasdb

import (
	"context"
	"github.com/dscamargo/rinha_backend_go/internal/domain/pessoa"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker struct {
	db            *pgxpool.Pool
	insertChannel chan pessoa.Pessoa
}

func NewInsertChannel() chan pessoa.Pessoa {
	return make(chan pessoa.Pessoa)
}

func NewWorker(db *pgxpool.Pool, insertChannel chan pessoa.Pessoa) *Worker {
	return &Worker{
		db:            db,
		insertChannel: insertChannel,
	}
}

func (w *Worker) Process() {
	for i := 0; i < 5; i++ {
		go func() {
			for p := range w.insertChannel {
				_, err := w.db.Exec(
					context.Background(),
					QueryInsertPessoa,
					p.ID,
					p.Apelido,
					p.Nome,
					p.Nascimento,
					p.StackStr(),
					p.Nome+" "+p.Apelido+" "+p.StackStr(),
				)
				if err != nil {
					log.Errorf("Erro ao inserir pessoa %v", err)
					return
				}
			}
		}()
	}
}
