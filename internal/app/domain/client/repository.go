package client

import "github.com/ryrden/rinha-de-backend-go/internal/app/data/models"

type Repository interface {
	FindByID(clientID string) (*Client, error)
	Update(client *Client) error
	CreateTransaction(clientID string, value int, kind string, description string) (*models.ClientTransactionResponse, error)
	GetClientExtract(clientID string) (*models.GetClientExtractResponse, error)
}
