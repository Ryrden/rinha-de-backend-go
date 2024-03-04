package client

import "github.com/ryrden/rinha-de-backend-go/internal/app/data/models"

type CreateTransaction struct {
	Repository Repository
}


func (c *CreateTransaction) Execute(clientID string, value int, kind string, description string) (*models.ClientTransactionResponse, error) {
	return c.Repository.CreateTransaction(clientID, value, kind, description)
}

func NewCreateTransaction(repository Repository) *CreateTransaction {
	return &CreateTransaction{Repository: repository}
}
