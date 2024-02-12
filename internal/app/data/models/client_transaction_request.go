package models

type CreateTransactionRequest struct {
	Value       int    `json:"valor" validate:"required,gt=0"`             // greater than 0
	Kind        string `json:"tipo" validate:"required,eq=c|eq=d"`         // c or d
	Description string `json:"descricao" validate:"required,gte=1,lte=10"` // 1 to 10 characters
}
