package main

import (
	"github.com/dscamargo/rinha_backend_go/config"
	"github.com/dscamargo/rinha_backend_go/internal/http"
	"github.com/dscamargo/rinha_backend_go/internal/http/controllers"
	"github.com/dscamargo/rinha_backend_go/internal/http/routers"
	"github.com/dscamargo/rinha_backend_go/internal/infra/db"
	"github.com/dscamargo/rinha_backend_go/pessoa"
	"github.com/dscamargo/rinha_backend_go/pessoa/postgres"
	"github.com/dscamargo/rinha_backend_go/pessoa/redis"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.Module,
		controllers.Module,
		routers.Module,
		http.Module,
		db.Module,
		postgres.Module,
		redis.Module,
		pessoa.Module,
		fx.Invoke(func(dispatcher *postgres.Dispatcher) {
			go dispatcher.Run()
		}),
		fx.Invoke(func(*fasthttp.Server) {}),
	)

	app.Run()
}
