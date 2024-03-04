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
}

func (c *ClientRepository) CreateTransaction(clientID string, value int, kind string, description string) (*models.ClientTransactionResponse, error) {
	log.Infof("Creating transaction for client %s", clientID)

	// Start a transaction
	log.Info("Starting transaction")
	tx, err := c.db.Begin(context.Background())
	if err != nil {
		log.Errorf("Error starting transaction: %s", err)
		return nil, err
	}
	defer tx.Rollback(context.Background())

	// Fetch client from database
	log.Infof("Attempting to find client by ID: %s", clientID)
	var clientResult client.Client
	log.Infof("Fetching client from database: %s", clientID)
	err = tx.QueryRow(
		context.Background(),
		PessimistLockingClientQuery,
		clientID,
	).Scan(&clientResult.ID, &clientResult.Limit, &clientResult.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, client.ErrClientNotFound
		}

		log.Errorf("Error finding client by ID: %s, error: %s", clientID, err)
		return nil, err
	}

	// Check if client can afford the transaction
	log.Infof("Checking if client can afford the transaction: %s, value: %d, kind: %s", clientID, value, kind)
	if !clientResult.CanAfford(value, kind) {
		log.Infof("Client cannot afford the transaction: %s, value: %d, kind: %s", clientID, value, kind)
		return nil, client.ErrClientCannotAfford
	}

	// Update client balance
	clientResult.AddTransaction(value, kind)
	log.Infof("Updating client balance: %s", clientID)
	_, err = tx.Exec(
		context.Background(),
		UpdateClientBalanceQuery,
		clientResult.Balance,
		clientResult.ID,
	)
	if err != nil {
		log.Errorf("Error updating client: %s, error: %s", clientID, err)
		return nil, err
	}

	// Insert a transcation record
	log.Infof("Inserting transaction record for client: %s, value: %d, kind: %s", clientID, value, kind)
	_, err = tx.Exec(
		context.Background(),
		InsertTransactionQuery,
		clientID,
		value,
		kind,
		description,
	)
	if err != nil {
		log.Errorf("Error creating transaction for client ID: %s, error: %s", clientID, err)
		return nil, err
	}

	log.Info("Committing transaction")
	err = tx.Commit(context.Background())
	if err != nil {
		log.Errorf("Error committing transaction: %s", err)
		return nil, err
	}

	log.Infof("Transaction successfully created for client ID: %s, value: %d, kind: %s", clientID, value, kind)
	return &models.ClientTransactionResponse{
		Limit:   clientResult.Limit,
		Balance: clientResult.Balance,
	}, nil
}

func (c *ClientRepository) GetClientExtract(clientID string) (*models.GetClientExtractResponse, error) {
	log.Infof("Starting to get client extract for clientID: %s", clientID)

	var transactions = new(models.GetClientExtractResponse)

	var clientResult client.Client
	log.Infof("Fetching client from database: %s", clientID)
	err := c.db.QueryRow(
		context.Background(),
		GetClientQuery,
		clientID,
	).Scan(&clientResult.ID, &clientResult.Limit, &clientResult.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, client.ErrClientNotFound
		}

		log.Errorf("Error finding client by ID: %s, error: %s", clientID, err)
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
		GetClientExtractQuery,
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

func NewClientRepository(db *pgxpool.Pool) client.Repository {
	return &ClientRepository{
		db: db,
	}
}
