package clientdb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3/log"
	"github.com/redis/rueidis"
	"github.com/ryrden/rinha-de-backend-go/internal/app/domain/client"
)

var ctx = context.Background()

type Cache struct {
	client rueidis.Client
}

func (p *Cache) Get(key string) (*client.Client, error) {
	getCmd := p.client.
		B().
		Get().
		Key("client:" + key).
		Build()

	clientBytes, err := p.client.Do(ctx, getCmd).AsBytes()
	if err != nil {
		return nil, err
	}

	var client client.Client
	err = sonic.Unmarshal(clientBytes, &client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (p *Cache) Set(client *client.Client) error {
	item, err := sonic.MarshalString(client)
	if err != nil {
		return err
	}

	setclientCmd := p.client.
		B().
		Set().
		Key("client:" + client.ID).
		Value(item).
		Ex(time.Minute).
		Build()

	cmds := make(rueidis.Commands, 0, 2)
	cmds = append(cmds, setclientCmd)

	for _, res := range p.client.DoMulti(ctx, cmds...) {
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
		panic(err)
	}

	log.Infof("Redis client created, connected to %s:%s", os.Getenv("CACHE_HOST"), os.Getenv("CACHE_PORT"))
	return &Cache{client: client}
}
