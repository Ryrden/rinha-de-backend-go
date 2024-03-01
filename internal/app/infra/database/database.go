package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/config"
)

var (
	db   *pgxpool.Pool
	once sync.Once
)

func NewPostgresDatabase(config *config.Config) *pgxpool.Pool {
	once.Do(func() {
		connectionUrl := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)

		poolConfig, err := pgxpool.ParseConfig(connectionUrl)
		if err != nil {
			log.Fatalf("Error parsing the connection URL: %s\n", err)
		}

		if config.Database.Max_db_connections != "" && config.Database.Min_db_connections != "" {
			if maxConnections, err := strconv.ParseInt(config.Database.Max_db_connections, 10, 64); err == nil {
				poolConfig.MaxConns = int32(maxConnections)
			}
			if minConnections, err := strconv.ParseInt(config.Database.Min_db_connections, 10, 64); err == nil {
				poolConfig.MinConns = int32(minConnections)
			}
		}

		db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatalf("Error connecting to the database: %s\n", err)
		}

		if err := db.Ping(context.Background()); err != nil {
			log.Fatalf("Error pinging the database: %s\n", err)
		}
	})

	return db

}
