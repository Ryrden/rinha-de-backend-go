package clientdb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3/log"
	"github.com/redis/rueidis"
	"github.com/ryrden/rinha-de-backend-go/internal/app/data/models"
	"github.com/ryrden/rinha-de-backend-go/internal/app/domain/client"
)

var ctx = context.Background()

type Cache struct {
	extract rueidis.Client
}

func (c *Cache) GetExtract(key string) (*models.GetClientExtractResponse, error) {
	getCmd := c.extract.
		B().
		Get().
		Key("clientExtract:" + key).
		Build()

	clientExtractBytes, err := c.extract.Do(ctx, getCmd).AsBytes()
	if err != nil {
		return nil, err
	}

	var clientExtract models.GetClientExtractResponse
	err = sonic.Unmarshal(clientExtractBytes, &clientExtract)
	if err != nil {
		return nil, err
	}

	return &clientExtract, nil
}

func (p *Cache) SetExtract(client *client.Client) error {
	item, err := sonic.MarshalString(client)
	if err != nil {
		return err
	}

	setclientCmd := p.extract.
		B().
		Set().
		Key("clientExtract:" + client.ID).
		Value(item).
		Ex(5 * time.Second).
		Build()

	cmds := make(rueidis.Commands, 0, 2)
	cmds = append(cmds, setclientCmd)

	for _, res := range p.extract.DoMulti(ctx, cmds...) {
		err := res.Error()

		if err != nil {
			return err
		}
	}

	return nil
}

func NewCache() *Cache {
	address := fmt.Sprintf(
		"%s:%s",
		os.Getenv("CACHE_HOST"),
		os.Getenv("CACHE_PORT"),
	)

	opts := rueidis.ClientOption{
		InitAddress:      []string{address},
		AlwaysPipelining: true,
	}
	client, err := rueidis.NewClient(opts)
	if err != nil {
		log.Errorf("Failed to create Redis client: %v", err)
		
	}

	log.Infof("Redis client created, connected to %s:%s", os.Getenv("CACHE_HOST"), os.Getenv("CACHE_PORT"))
	return &Cache{extract: client}
}
