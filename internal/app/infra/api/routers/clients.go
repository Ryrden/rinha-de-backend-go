package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/api/controllers"
)

type ClientRouter struct {
	controller *controllers.ClientController
}

func (c *ClientRouter) Load(r *fiber.App) {
	r.Get("/clientes/:id/extrato", c.controller.GetClientTransactions)
	r.Post("/clientes/:id/transacoes", c.controller.CreateTransaction)
}

func NewClientRouter(
	controller *controllers.ClientController,
) *ClientRouter {
	return &ClientRouter{
		controller: controller,
	}
}
