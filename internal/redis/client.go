package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func New() *redis.Client {
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}
	route := os.Getenv("REDIS_ROUTE")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", route, port),
		Password: "",
		DB:       0,
	})
	var ctx = context.Background()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	return client
}
