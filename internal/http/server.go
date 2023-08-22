package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"os"
	"runtime/pprof"
)

func prof() func() {
	f, err := os.Create(os.Getenv("CPU_PROF"))
	if err != nil {
		log.Fatal(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Info("Error starting CPU profile: %v", err)
	}

	mf, err := os.Create(os.Getenv("MEM_PROF"))
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(mf)

	return func() {
		pprof.StopCPUProfile()
		f.Close()
		mf.Close()
	}
}

func NewHTTPServer(lifecycle fx.Lifecycle, router *fiber.App) *fasthttp.Server {
	//var shutdown func()

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {

				//shutdown = prof()

				log.Info("Starting web server...")
				if err := router.Listen(":" + os.Getenv("PORT")); err != nil {
					log.Fatalf("Error starting the server: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			//defer shutdown()
			log.Info("Stopping web server...")
			return router.ShutdownWithContext(ctx)
		},
	})

	return router.Server()
}
