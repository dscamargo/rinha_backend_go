package routers

import (
	"github.com/gofiber/fiber/v2"
	"rinha_v2/internal/http/controllers"
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
