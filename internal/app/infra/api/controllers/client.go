package controllers

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/ryrden/rinha-de-backend-go/internal/app/data/models"
	"github.com/ryrden/rinha-de-backend-go/internal/app/domain/client"
)

type ClientController struct {
	createTransaction *client.CreateTransaction
	getClientExtract  *client.GetClientExtract
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

	clientResult, err := c.createTransaction.Execute(id, dto.Value, dto.Kind, dto.Description)
	if err != nil {
		if errors.Is(err, client.ErrClientNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if errors.Is(err, client.ErrClientCannotAfford) {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(clientResult)

}

func (c *ClientController) GetClientExtract(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	clientExtract, err := c.getClientExtract.Execute(id)
	if err != nil {
		if errors.Is(err, client.ErrClientNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(clientExtract)
}

func NewClientController(
	createTransaction *client.CreateTransaction,
	getClientExtract *client.GetClientExtract,
) *ClientController {
	return &ClientController{
		createTransaction: createTransaction,
		getClientExtract:  getClientExtract,
	}
}
