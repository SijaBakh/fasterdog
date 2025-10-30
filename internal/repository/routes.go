package repository

import (
	"context"

	"github.com/SijaBakh/fasterdog/internal/models"
)

func (fr *FasterdogRepository) GetRoutes(ctx context.Context) ([]models.Route, error) {
	routes, err := fr.db.GetRoutes(ctx)
	if err != nil {
		return nil, err
	}

	return routes, nil
}

func (fr *FasterdogRepository) ExecuteManyRoutes(ctx context.Context, routes []models.Route) error {
	err := fr.db.ExecuteManyRoutes(ctx, routes)
	if err != nil {
		return err
	}

	return nil
}
