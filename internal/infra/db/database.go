package db

import (
	"context"
	"github.com/dscamargo/rinha_backend_go/config"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

var (
	db   *pgxpool.Pool
	once sync.Once
)

func NewPGDatabase(cfg *config.Config) *pgxpool.Pool {
	once.Do(func() {
		dbUrl := cfg.Database.Url
		poolConfig, err := pgxpool.ParseConfig(dbUrl)
		if err != nil {
			log.Fatal("Erro ao buscar configurações poolConfig")
		}

		db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatal("Erro ao conectar no banco de dados")
		}
	})

	return db
}
