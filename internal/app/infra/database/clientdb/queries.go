package clientdb

const (
	PessimistLockingClientQuery = "SELECT * FROM clients WHERE id = $1 FOR UPDATE" // Pessimist locking with "FOR UPDATE"
	UpdateClientBalanceQuery = "UPDATE clients SET balance = $1 WHERE id = $2"
	InsertTransactionQuery = "INSERT INTO transactions (client_id, value, description) VALUES ($1, $2, $3)"

	GetClientExtractQuery = "SELECT id, balance_limit, balance FROM clients WHERE id = $1"
)