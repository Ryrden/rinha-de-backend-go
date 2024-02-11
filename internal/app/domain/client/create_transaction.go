package client

type CreateTransaction struct {
	Repository Repository
}

// TODO: A lógica vem aqui

func (c *CreateTransaction) Execute(clientID string, value int, kind string, description string) (*Client, error) {
	return c.Repository.CreateTransaction(clientID, value, kind, description)
}

func NewCreateTransaction(repository Repository) *CreateTransaction {
	return &CreateTransaction{Repository: repository}
}
