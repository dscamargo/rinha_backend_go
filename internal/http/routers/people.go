package routers

import (
	"github.com/dscamargo/rinha_backend_go/internal/http/controllers"
	"github.com/gofiber/fiber/v2"
)

type PessoaRouter struct {
	controller *controllers.PessoaController
}

func (ro *PessoaRouter) Load(r *fiber.App) {
	r.Get("/pessoas", ro.controller.List)
	r.Get("/pessoas/:id", ro.controller.Show)
	r.Post("/pessoas", ro.controller.Create)
	r.Get("/contagem-pessoas", ro.controller.Count)
}

func NewPessoaRouter(controller *controllers.PessoaController) *PessoaRouter {
	return &PessoaRouter{
		controller,
	}
}
