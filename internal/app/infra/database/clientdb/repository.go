package clientdb

import (
	"context"
	"errors"

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
	var client client.Client

	err := c.db.QueryRow(
		context.Background(),
		"SELECT id, balance_limit, balance FROM clients WHERE id = $1",
		id,
	).Scan(&client.ID, &client.Limit, &client.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("client not found")
		}

		log.Errorf("Error finding client by id: %s", err)
		return nil, err
	}

	return &client, nil
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
	client, err := c.FindByID(clientID)
	if err != nil {
		return nil, err
	}
	client.AddTransaction(value, kind)
	err = c.Update(client)
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
		Limit:   client.Limit,
		Balance: client.Balance,
	}, nil
}

func (c *ClientRepository) GetClientExtract(clientID string) (*models.GetClientExtractResponse, error) {
	var transactions *models.GetClientExtractResponse

	rows, err := c.db.Query(
		context.Background(),
		"SELECT id, value, kind, description FROM transactions WHERE client_id = $1",
		clientID,
	)
	if err != nil {
		log.Errorf("Error getting client transactions: %s", err)
		return nil, err
	}
	defer rows.Close()

	// for rows.Next() {
	// 	var transaction models.Transaction

	// 	err := rows.Scan(&transaction.ID, &transaction.Value, &transaction.Kind, &transaction.Description)
	// 	if err != nil {
	// 		log.Errorf("Error scanning transaction: %s", err)
	// 		return nil, err
	// 	}

	// 	transactions = append(transactions, &transaction)
	// }

	return transactions, nil
}

func NewClientRepository(db *pgxpool.Pool) client.Repository {
	return &ClientRepository{db: db}
}
