package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"os"
	"rinha_v2/internal/domain/pessoa"
	"rinha_v2/internal/http"
	"rinha_v2/internal/http/controllers"
	"rinha_v2/internal/http/routers"
	"rinha_v2/internal/infra/db"
	"rinha_v2/internal/infra/db/pessoasdb"
)

func main() {
	fmt.Printf("Process id of %d main \n", os.Getpid())

	app := fx.New(
		controllers.Module,
		routers.Module,
		http.Module,
		db.Module,
		pessoa.Module,
		fx.Invoke(func(c *pessoasdb.Worker) {
			go c.Process()
		}),
		fx.Invoke(func(*fasthttp.Server) {}),
	)

	app.Run()
}
