package routes

import (
	"context"
	"net/http"
	"slices"

	"github.com/SijaBakh/fasterdog/internal/adapter/db"
	"github.com/SijaBakh/fasterdog/internal/models"
	"github.com/SijaBakh/fasterdog/internal/repository"

	"github.com/go-chi/chi/v5"
)

type Route interface {
	Method() string
	Path() string
}

func NewRoute(m, p string) Route {
	return models.NewRoute(m, p)
}

func GetRoutes(r chi.Routes) ([]models.Route, error) {
	var routes []models.Route
	walkFunc := func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routes = append(routes, *models.NewRoute(method, route))
		return nil
	}

	err := chi.Walk(r, walkFunc)
	return routes, err
}

func CheckRoutes(dsn string, routes []models.Route) error {
	ctx := context.Background()
	db, err := db.New(dsn, ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	fr := repository.New(db)
	dbRoutes, err := fr.GetRoutes(ctx)
	if err != nil {
		return err
	}

	if len(dbRoutes) == 0 {
		err := fr.ExecuteManyRoutes(ctx, routes)
		return err
	}

	difRoutes := difference(routes, dbRoutes)
	if len(difRoutes) > 0 {
		err := fr.ExecuteManyRoutes(ctx, difRoutes)
		return err
	}

	return nil
}

func difference(routes, dbRoutes []models.Route) []models.Route {
	difRoutes := make([]models.Route, len(routes))
	for _, r := range routes {
		if slices.Contains(dbRoutes, r) {
			continue
		}
		difRoutes = append(difRoutes, r)
	}
	return difRoutes
}
