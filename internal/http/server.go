package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"os"
)

func NewHTTPServer(lifecycle fx.Lifecycle, router *fiber.App) *fasthttp.Server {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info("Starting web server...")
				if err := router.Listen(":" + os.Getenv("PORT")); err != nil {
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
