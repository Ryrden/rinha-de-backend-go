package clientdb

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ryrden/rinha-de-backend-go/internal/app/data/models"
	"github.com/ryrden/rinha-de-backend-go/internal/app/domain/client"
)

type ClientRepository struct {
	db *pgxpool.Pool
	// TODO: add cache and jobQueue
}

func (c *ClientRepository) FindByID(id string) (*client.Client, error) {
	var clientResult client.Client

	err := c.db.QueryRow(
		context.Background(),
		"SELECT id, balance_limit, balance FROM clients WHERE id = $1",
		id,
	).Scan(&clientResult.ID, &clientResult.Limit, &clientResult.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, client.ErrClientNotFound
		}

		log.Errorf("Error finding client by id: %s", err)
		return nil, err
	}

	return &clientResult, nil
}

func (c *ClientRepository) Update(client *client.Client) error {
	_, err := c.db.Exec(
		context.Background(),
		"UPDATE clients SET balance = $1 WHERE id = $2",
		client.Balance,
		client.ID,
	)
	if err != nil {
		log.Errorf("Error updating client: %s", err)
		return err
	}

	return nil
}

func (c *ClientRepository) CreateTransaction(clientID string, value int, kind string, description string) (*models.ClientTransactionResponse, error) {
	clientResult, err := c.FindByID(clientID)
	if err != nil {
		return nil, err
	}
	if clientResult == nil {
		return nil, client.ErrClientNotFound
	}

	if !clientResult.CanAfford(value, kind) {
		return nil, client.ErrClientCannotAfford
	}

	clientResult.AddTransaction(value, kind)
	err = c.Update(clientResult)
	if err != nil {
		return nil, err
	}

	_, err = c.db.Exec(
		context.Background(),
		"INSERT INTO transactions(client_id, amount, kind, description) VALUES($1, $2, $3, $4)",
		clientID,
		value,
		kind,
		description,
	)
	if err != nil {
		log.Errorf("Error creating transaction: %s", err)
		return nil, err
	}

	return &models.ClientTransactionResponse{
		Limit:   clientResult.Limit,
		Balance: clientResult.Balance,
	}, nil
}

func (c *ClientRepository) GetClientExtract(clientID string) (*models.GetClientExtractResponse, error) {
	var transactions = new(models.GetClientExtractResponse)

	clientResult, err := c.FindByID(clientID)
	if err != nil {
		return nil, err
	}
	if clientResult == nil {
		return nil, client.ErrClientNotFound
	}
	transactions.Balance = models.Balance{
		Total:       clientResult.Balance,
		ExtractedAt: time.Now().Format("2006-01-02T15:04:05.000000Z"),
		Limit:       clientResult.Limit,
	}

	rows, err := c.db.Query(
		context.Background(),
		"SELECT amount, kind, description, created_at FROM transactions WHERE client_id = $1 ORDER BY created_at DESC LIMIT 10",
		clientID,
	)

	if err != nil {
		log.Errorf("Error getting client extract: %s", err)
		return nil, err
	}
	defer rows.Close()

	transactions.Transactions = make([]models.Transaction, 0)
	for rows.Next() {
		var amount int
		var kind string
		var description string
		var createdAt time.Time

		if err := rows.Scan(&amount, &kind, &description, &createdAt); err != nil {
			log.Errorf("Error scanning transaction: %s", err)
			return nil, err
		}

		var transaction models.Transaction
		transaction.Value = amount
		transaction.Kind = kind
		transaction.Description = description
		transaction.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000000Z")

		transactions.Transactions = append(transactions.Transactions, transaction)
	}

	return transactions, nil
}

func NewClientRepository(db *pgxpool.Pool) client.Repository {
	return &ClientRepository{db: db}
}
