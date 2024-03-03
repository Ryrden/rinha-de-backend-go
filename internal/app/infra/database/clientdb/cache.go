package clientdb

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3/log"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	conn *redis.Client
}

func NewCache() *Cache {
	address := fmt.Sprintf(
		"%s:%s",
		os.Getenv("CACHE_HOST"),
		os.Getenv("CACHE_PORT"),
	)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	log.Infof("Redis client created, connected to %s:%s", os.Getenv("CACHE_HOST"), os.Getenv("CACHE_PORT"))
	return &Cache{conn: client}
}
