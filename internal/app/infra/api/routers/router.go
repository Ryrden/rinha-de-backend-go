package routers

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/config"
)

type Router interface {
	Load()
}

func MakeRouter(
	clientRouter *ClientRouter,
	config *config.Config,
) *fiber.App {
	cfg := fiber.Config{
		AppName:       "rinha-go by @ryrden",
		CaseSensitive: true,
		// FIXME: This is not working
		//Prefork:       config.Server.Prefork,
	}

	if config.Server.UseSonic {
		log.Info("Loading Sonic JSON into the router")
		cfg.JSONEncoder = sonic.Marshal
		cfg.JSONDecoder = sonic.Unmarshal
	}

	r := fiber.New(cfg)

	clientRouter.Load(r)

	return r
}
