package models

type CreateTransactionRequest struct {
	Value       *int    `json:"valor"`
	Kind        *string `json:"tipo"`
	Description *string `json:"descricao"`
}
