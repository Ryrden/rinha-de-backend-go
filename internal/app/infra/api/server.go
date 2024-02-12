package api

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/config"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func NewServer(
	lifecycle fx.Lifecycle,
	router *fiber.App,
	config *config.Config,
) *fasthttp.Server {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info("Starting the server...")

				// TODO: Profilling

				addr := fmt.Sprintf(":%s", config.Server.Port)
				if err := router.Listen(addr); err != nil {
					log.Fatalf("Error starting the server: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping the server...")
			return router.ShutdownWithContext(ctx)
		},
	})
	return router.Server()
}
