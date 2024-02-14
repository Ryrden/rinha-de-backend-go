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
	db       *pgxpool.Pool
	cache    *Cache
	jobQueue JobQueue
}

func (c *ClientRepository) FindByID(id string) (*client.Client, error) {
	log.Infof("Attempting to find client by ID: %s", id)

	var clientResult client.Client
	cachedClient, err := c.cache.GetClient(id)
	if err == nil {
		log.Infof("Client found in cache: %s", id)
		return cachedClient, nil
	}

	log.Infof("Fetching client from database: %s", id)
	err = c.db.QueryRow(
		context.Background(),
		"SELECT id, balance_limit, balance FROM clients WHERE id = $1",
		id,
	).Scan(&clientResult.ID, &clientResult.Limit, &clientResult.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, client.ErrClientNotFound
		}

		log.Errorf("Error finding client by ID: %s, error: %s", id, err)
		return nil, err
	}

	c.cache.SetClient(&clientResult)

	log.Infof("Client successfully found and returned: %s", id)
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

	c.cache.SetClient(client)

	return nil
}

func (c *ClientRepository) CreateTransaction(clientID string, value int, kind string, description string) (*models.ClientTransactionResponse, error) {
	log.Infof("Creating transaction for client %s", clientID)
	clientResult, err := c.FindByID(clientID)
	if err != nil {
		return nil, err
	}

	if !clientResult.CanAfford(value, kind) {
		log.Infof("Client cannot afford the transaction: %s, value: %d, kind: %s", clientID, value, kind)
		return nil, client.ErrClientCannotAfford
	}

	clientResult.AddTransaction(value, kind)
	err = c.Update(clientResult)
	if err != nil {
		log.Errorf("Error updating client: %s, error: %s", clientID, err)
		return nil, err
	}

	job := Job{
		Payload: &ClientTransactionPayload{
			Client:      clientResult,
			Value:       value,
			Kind:        kind,
			Description: description,
		},
	}

	c.jobQueue <- job

	log.Infof("Transaction successfully created for client ID: %s, value: %d, kind: %s", clientID, value, kind)
	return &models.ClientTransactionResponse{
		Limit:   clientResult.Limit,
		Balance: clientResult.Balance,
	}, nil
}

func (c *ClientRepository) GetClientExtract(clientID string) (*models.GetClientExtractResponse, error) {
	log.Infof("Starting to get client extract for clientID: %s", clientID)

	var transactions = new(models.GetClientExtractResponse)

	clientResult, err := c.FindByID(clientID)
	if err != nil {
		return nil, err
	}

	transactions.Balance = models.Balance{
		Total:       clientResult.Balance,
		ExtractedAt: time.Now().Format(time.RFC3339),
		Limit:       clientResult.Limit,
	}

	log.Infof("Fetching last 10 transactions for clientID: %s", clientID)
	rows, err := c.db.Query(
		context.Background(),
		"SELECT amount, kind, description, created_at FROM transactions WHERE client_id = $1 ORDER BY created_at DESC LIMIT 10",
		clientID,
	)
	if err != nil {
		log.Errorf("Error getting client extract for clientID: %s, error: %s", clientID, err)
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
			log.Errorf("Error scanning transaction for clientID: %s, error: %s", clientID, err)
			return nil, err
		}

		transaction := models.Transaction{
			Value:       amount,
			Kind:        kind,
			Description: description,
			CreatedAt:   createdAt.Format("2006-01-02T15:04:05.000000Z"),
		}

		transactions.Transactions = append(transactions.Transactions, transaction)
	}

	log.Infof("Successfully retrieved client extract for clientID: %s", clientID)
	return transactions, nil
}

func NewClientRepository(db *pgxpool.Pool, jobQueue JobQueue, cache *Cache) client.Repository {
	return &ClientRepository{
		db:       db,
		jobQueue: jobQueue,
		cache:    cache,
	}
}
