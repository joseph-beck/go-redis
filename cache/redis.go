package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Lifetime int
}

type Cache struct {
	Client *redis.Client
	Config Config
}

func New() *Cache {
	c := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &Cache{
		Client: c,
	}
}

func (c *Cache) Ping() (string, error) {
	ctx := context.Background()
	s, err := c.Client.Ping(ctx).Result()
	if err != nil {
		return "", err
	}

	return s, nil
}

func (c *Cache) Close() error {
	fmt.Println("closing redis connection")
	err := c.Client.Close()
	return err
}

func (c *Cache) Set() {

}

func (c *Cache) Get() {

}
