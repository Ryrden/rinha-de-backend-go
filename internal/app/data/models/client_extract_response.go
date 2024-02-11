package models

type GetClientExtractResponse struct {
	Balance      Balance       `json:"saldo"`
	Transactions []Transaction `json:"ultimas_transacoes"`
}

type Balance struct {
	Total       int    `json:"total"`
	ExtractedAt string `json:"data_extrato"`
	Limit       int    `json:"limite"`
}

type Transaction struct {
	Value       int    `json:"valor"`
	Kind        string `json:"tipo"`
	Description string `json:"descricao"`
	PerformedAt string `json:"realizada_em"`
}
