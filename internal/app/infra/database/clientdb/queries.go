package clientdb

const (
	PessimistLockingClientQuery = "SELECT id, balance_limit, balance FROM clients WHERE id = $1 FOR UPDATE" // Pessimist locking with "FOR UPDATE"
	UpdateClientBalanceQuery    = "UPDATE clients SET balance = $1 WHERE id = $2"
	InsertTransactionQuery      = "INSERT INTO transactions(client_id, amount, kind, description) VALUES($1, $2, $3, $4)"
	GetClientQuery              = "SELECT id, balance_limit, balance FROM clients WHERE id = $1"
	GetClientExtractQuery       = "SELECT amount, kind, description, created_at FROM transactions WHERE client_id = $1 ORDER BY created_at DESC LIMIT 10"
)
