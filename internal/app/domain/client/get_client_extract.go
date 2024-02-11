package client

import "github.com/ryrden/rinha-de-backend-go/internal/app/data/models"

type GetClientExtract struct {
	Repository Repository
}

func (c *GetClientExtract) Execute(clientID string) (*models.GetClientExtractResponse, error) {
	return c.Repository.GetClientExtract(clientID)
}

func NewGetClientExtract(repository Repository) *GetClientExtract {
	return &GetClientExtract{Repository: repository}
}
