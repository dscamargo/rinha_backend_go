package routers

import (
	"github.com/gofiber/fiber/v2"
)

func MakeRouter(
	pessoaRouter *PessoaRouter) *fiber.App {
	cfg := fiber.Config{AppName: "rinha_backend"}
	router := fiber.New(cfg)
	//router.Use(logger.New(logger.Config{
	//	Format: "${pid} - ${latency} - ${locals:requestid} ${status} - Query: ${queryParams} - ${method} ${path}​\n",
	//}))

	pessoaRouter.Load(router)

	return router

}
