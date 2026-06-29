package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	startupCtx := context.Background()

	_, err := rdb.Ping(startupCtx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	fmt.Println("Connected to Redis successfully!")
	return rdb, nil
}
