package routers

import (
	"github.com/bytedance/sonic"
	"github.com/dscamargo/rinha_backend_go/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func MakeRouter(pessoaRouter *PessoaRouter, cfg *config.Config) *fiber.App {
	fiberCfg := fiber.Config{AppName: "rinha_backend"}

	if cfg.App.EnableSonicJson {
		log.Info("[MakeRouter] - Loading Sonic JSON...")
		fiberCfg.JSONEncoder = sonic.Marshal
		fiberCfg.JSONDecoder = sonic.Unmarshal
	}

	router := fiber.New(fiberCfg)

	pessoaRouter.Load(router)

	return router

}
