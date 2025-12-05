package redis

import (
	"context"

	"github.com/SijaBakh/fasterdog/internal/models"

	"github.com/redis/go-redis/v9"
)

type RedisClientInterfaces interface {
	RGetPermissions(ctx context.Context, username string) (*models.PermissionsResult, error)
}

type RedisClient struct {
	client *redis.Client
}

func New(url string, max_pool int) (RedisClientInterfaces, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	opts.MaxActiveConns = max_pool
	return &RedisClient{
		client: redis.NewClient(opts),
	}, nil
}
