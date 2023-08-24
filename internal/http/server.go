package http

import (
	"context"
	"github.com/dscamargo/rinha_backend_go/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func NewHTTPServer(lifecycle fx.Lifecycle, router *fiber.App, cfg *config.Config) *fasthttp.Server {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info("Starting web server...")
				if err := router.Listen(":" + cfg.App.Port); err != nil {
					log.Fatalf("Error starting the server: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping web server...")
			return router.ShutdownWithContext(ctx)
		},
	})

	return router.Server()
}
