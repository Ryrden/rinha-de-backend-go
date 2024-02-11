package client

import "github.com/ryrden/rinha-de-backend-go/internal/app/data/models"

type GetClientTransactions struct {
	Repository Repository
}

func (c *GetClientTransactions) Execute(clientID string) ([]models.Transaction, error) {
	return c.Repository.GetClientTransactions(clientID)
}

func NewGetClientTransactions(repository Repository) *GetClientTransactions {
	return &GetClientTransactions{Repository: repository}
}
