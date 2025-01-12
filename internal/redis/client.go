package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(lc fx.Lifecycle) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Connecting to Redis...")
			if err := client.Ping(ctx).Err(); err != nil {
				return fmt.Errorf("failed to connect to Redis: %w", err)
			}
			fmt.Println("Connected to Redis successfully.")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Closing Redis connection...")
			return client.Close()
		},
	})

	return &RedisClient{Client: client}
}
