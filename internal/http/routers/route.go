package routers

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"os"
)

func MakeRouter(
	pessoaRouter *PessoaRouter) *fiber.App {
	cfg := fiber.Config{AppName: "rinha_backend"}

	if os.Getenv("ENABLE_SONIC_JSON") == "1" {
		log.Info("[MakeRouter] - Loading Sonic JSON...")
		cfg.JSONEncoder = sonic.Marshal
		cfg.JSONDecoder = sonic.Unmarshal
	}

	router := fiber.New(cfg)

	pessoaRouter.Load(router)

	return router

}
