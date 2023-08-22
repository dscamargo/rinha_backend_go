package main

import (
	"fmt"
	"github.com/dscamargo/rinha_backend_go/internal/domain/pessoa"
	"github.com/dscamargo/rinha_backend_go/internal/http"
	"github.com/dscamargo/rinha_backend_go/internal/http/controllers"
	"github.com/dscamargo/rinha_backend_go/internal/http/routers"
	"github.com/dscamargo/rinha_backend_go/internal/infra/db"
	"github.com/dscamargo/rinha_backend_go/internal/infra/db/pessoasdb"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"os"
)

func main() {
	fmt.Printf("Process id of %d main \n", os.Getpid())

	app := fx.New(
		controllers.Module,
		routers.Module,
		http.Module,
		db.Module,
		pessoa.Module,
		fx.Invoke(func(dispatcher *pessoasdb.Dispatcher) {
			go dispatcher.Run()
		}),
		fx.Invoke(func(*fasthttp.Server) {}),
	)

	app.Run()
}
