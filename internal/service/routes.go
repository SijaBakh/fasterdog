package service

import (
	"context"

	"github.com/SijaBakh/fasterdog/internal/models"
)

func (fs *FasterdogService) GetRoutes(ctx context.Context) ([]models.Route, error) {
	routes, err := fs.repo.GetRoutes(ctx)
	if err != nil {
		return nil, err
	}

	return routes, nil
}

func (fs *FasterdogService) ExecuteManyRoutes(ctx context.Context, routes []models.Route) error {
	err := fs.repo.ExecuteManyRoutes(ctx, routes)
	if err != nil {
		return err
	}

	return nil
}
