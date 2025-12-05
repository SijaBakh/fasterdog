package service

import (
	"context"

	"github.com/SijaBakh/fasterdog/internal/infrastructure/db"
	"github.com/SijaBakh/fasterdog/internal/infrastructure/redis"
	"github.com/SijaBakh/fasterdog/internal/models"
	"github.com/SijaBakh/fasterdog/internal/repository"
)

type FasterdogServiceInterfaces interface {
	GetPermissions(ctx context.Context, userName, domainName string) (*models.PermissionsResult, error)
	GetRoutes(ctx context.Context) ([]models.Route, error)
	ExecuteManyRoutes(ctx context.Context, routes []models.Route) error
	RGetPermissions(ctx context.Context, username string) (*models.PermissionsResult, error)
}

type FasterdogService struct {
	repo    repository.FasterdogRepositoryInterfaces
	rClient redis.RedisClientInterfaces
}

func New(ctx context.Context, redisDSN, authDSN string, redisMP int) (FasterdogServiceInterfaces, error) {
	rc, err := redis.New(redisDSN, redisMP)
	if err != nil {
		return nil, err
	}

	db, err := db.New(authDSN, ctx)
	if err != nil {
		return nil, err
	}
	fr := repository.New(db)

	return &FasterdogService{
		repo:    fr,
		rClient: rc,
	}, nil
}
