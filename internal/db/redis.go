package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

var RedisClient *redis.Client

// ConnectRedis initializes a Redis client
func ConnectRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379" // fallback for local dev
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(fmt.Sprintf("failed to parse redis url: %v", err))
	}

	RedisClient = redis.NewClient(opt)

	// Test connection
	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	fmt.Println("✅ Connected to Redis")
}
