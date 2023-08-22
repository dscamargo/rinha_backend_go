package db

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"os"
	"sync"
)

var (
	db   *pgxpool.Pool
	once sync.Once
)

func NewPGDatabase() *pgxpool.Pool {
	once.Do(func() {
		dbUrl := os.Getenv("DATABASE_URL")
		poolConfig, err := pgxpool.ParseConfig(dbUrl)
		if err != nil {
			log.Fatal("Erro ao buscar configurações pgx")
		}

		conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatal("Erro ao conectar no banco de dados")
		}

		db = conn
	})

	return db
}
