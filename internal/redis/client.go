package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
)

type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type client struct {
	c *redis.Client
}

func NewClient(config config.RedisConfig) Client {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       int(config.Db),
	})

	return &client{c: c}
}

func (c *client) Get(ctx context.Context, key string) (string, error) {
	res, err := c.c.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}

	return res, err
}

func (c *client) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return c.c.Set(ctx, key, value, exp).Err()
}

func (c *client) Del(ctx context.Context, keys ...string) error {
	return c.c.Del(ctx, keys...).Err()
}
