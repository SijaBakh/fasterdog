package repository

import (
	"context"

	"github.com/SijaBakh/fasterdog/internal/infrastructure/db"
	"github.com/SijaBakh/fasterdog/internal/models"
)

type FasterdogRepositoryInterfaces interface {
	GetPermissions(ctx context.Context, userName, domainName string) ([]byte, error)
	GetRoutes(ctx context.Context) ([]models.Route, error)
	ExecuteManyRoutes(ctx context.Context, routes []models.Route) error
}

type FasterdogRepository struct {
	db *db.DB
}

func New(db *db.DB) FasterdogRepositoryInterfaces {
	return &FasterdogRepository{
		db: db,
	}
}
