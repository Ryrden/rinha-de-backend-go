package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/ryrden/rinha-de-backend-go/internal/app/data/models"
	"github.com/ryrden/rinha-de-backend-go/internal/app/domain/client"
)

type ClientController struct {
	createTransaction     *client.CreateTransaction
	getClientTransactions *client.GetClientTransactions
}

func (c *ClientController) CreateTransaction(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	var dto models.CreateTransactionRequest

	if err := ctx.Bind().Body(&dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	validator := validator.New()

	if err := validator.Struct(dto); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	client, err := c.createTransaction.Execute(id, dto.Value, dto.Kind, dto.Description)
	if err != nil {
		// TODO: treat specific errors
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(client)

}

func (c *ClientController) GetClientTransactions(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "not implemented",
	})
}

func NewClientController(
	createTransaction *client.CreateTransaction,
	getClientTransactions *client.GetClientTransactions,
) *ClientController {
	return &ClientController{
		createTransaction:     createTransaction,
		getClientTransactions: getClientTransactions,
	}
}
