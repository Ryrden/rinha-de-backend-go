package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/ryrden/rinha-de-backend-go/internal/app/domain/client"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/api"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/api/controllers"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/api/routers"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/config"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/database"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/database/clientdb"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func main() {
	uuid.EnableRandPool()

	if err := godotenv.Load(); err != nil {
		log.Warn("Could not load .env file")
	}

	fmt.Print("Starting Rinha de Backend\n")

	app := fx.New(
		config.Module,
		controllers.Module,
		routers.Module,
		api.Module,
		database.Module,
		client.Module,
		fx.Invoke(func(dispatcher *clientdb.Dispatcher) {
			go dispatcher.Run()
		}),
		fx.Invoke(func(*fasthttp.Server) {}),
		// fx.NopLogger,
	)

	app.Run()
}
