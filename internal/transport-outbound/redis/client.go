package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"os"
)

func Init() (*redis.Client, error) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		return nil, errors.New("REDIS_URL var not found")
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}
