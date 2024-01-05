package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

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
		Config: Config{Lifetime: 1},
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

func (c *Cache) Set(i interface{}, k string) (string, error) {
	ctx := context.Background()

	m, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	s, err := c.Client.Set(ctx, k, m, time.Duration(c.Config.Lifetime*int(time.Minute))).Result()
	if err != nil {
		return "", err
	}

	return s, nil
}

func (c *Cache) Get(i interface{}, k string) error {
	ctx := context.Background()
	r, err := c.Client.Get(ctx, k).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(r), &i)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Delete(k string) (int, error) {
	ctx := context.Background()
	s, err := c.Client.Del(ctx, k).Result()
	if err != nil {
		return 0, err
	}

	return int(s), nil
}

func (c *Cache) Contains(k string) (bool, error) {
	ctx := context.Background()
	e, err := c.Client.Exists(ctx, k).Result()
	if err != nil {
		return false, err
	}

	return e == 1, nil
}
