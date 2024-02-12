package controllers

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/ryrden/rinha-de-backend-go/internal/app/data/models"
	"github.com/ryrden/rinha-de-backend-go/internal/app/domain/client"
)

type ClientController struct {
	createTransaction *client.CreateTransaction
	getClientExtract  *client.GetClientExtract
}

func (c *ClientController) CreateTransaction(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	// verify if id is a integer
	// if id is not a integer, return a 400 status code
	if _, err := strconv.Atoi(id); err != nil {
		log.Warnf("Invalid client ID: %s", id)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Invalid client ID",
		})
	}
	log.Infof("Starting CreateTransaction for client ID: %s", id)

	var dto models.CreateTransactionRequest

	if err := ctx.Bind().Body(&dto); err != nil {
		log.Errorf("Error parsing request body for CreateTransaction, client ID: %s, error: %s", id, err)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	validator := validator.New()
	if err := validator.Struct(dto); err != nil {
		log.Warnf("Validation failed for CreateTransaction, client ID: %s, error: %s", id, err)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Validation failed",
		})
	}

	clientResult, err := c.createTransaction.Execute(id, dto.Value, dto.Kind, dto.Description)
	if err != nil {
		if errors.Is(err, client.ErrClientCannotAfford) {
			log.Warnf("Client cannot afford transaction, client ID: %s, error: %s", id, err)
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "Client cannot afford transaction",
			})
		}
		if errors.Is(err, client.ErrClientNotFound) {
			log.Warnf("Client not found, client ID: %s", id)
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Client not found",
			})
		}
		log.Errorf("Error executing CreateTransaction, client ID: %s, error: %s", id, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	log.Infof("Transaction created successfully for client ID: %s", id)
	return ctx.Status(fiber.StatusOK).JSON(clientResult)
}

func (c *ClientController) GetClientExtract(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		log.Warnf("Invalid client ID: %s", id)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Invalid client ID",
		})
	}
	log.Infof("Request to get client extract for client ID: %s", id)

	clientExtract, err := c.getClientExtract.Execute(id)
	if err != nil {
		if errors.Is(err, client.ErrClientNotFound) {
			log.Warnf("Client not found for ID: %s", id)
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Client not found",
			})
		}
		log.Errorf("Error retrieving client extract for ID: %s, error: %s", id, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	log.Infof("Successfully retrieved client extract for client ID: %s", id)
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
