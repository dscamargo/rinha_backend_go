package routers

import (
	"github.com/gofiber/fiber/v2"
)

func MakeRouter(
	pessoaRouter *PessoaRouter) *fiber.App {
	cfg := fiber.Config{AppName: "rinha_backend"}
	router := fiber.New(cfg)

	pessoaRouter.Load(router)

	return router

}
