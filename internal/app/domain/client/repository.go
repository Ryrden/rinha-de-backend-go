package client

import "github.com/ryrden/rinha-de-backend-go/internal/app/data/models"

type Repository interface {
	CreateTransaction(clientID string, value int, kind string, description string) (*Client, error)
	GetClientTransactions(clientID string) ([]models.Transaction, error)
}
