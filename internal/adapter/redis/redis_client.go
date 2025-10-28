package redis

import (
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func New(url string, max_pool int) (*RedisClient, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	opts.MaxActiveConns = max_pool
	return &RedisClient{
		client: redis.NewClient(opts),
	}, nil
}
