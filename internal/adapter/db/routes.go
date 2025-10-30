package db

import (
	"context"

	"github.com/SijaBakh/fasterdog/internal/models"

	"github.com/jackc/pgx/v5"
)

func (db *DB) GetRoutes(ctx context.Context) ([]models.Route, error) {
	query := `
	SELECT
        method,
		path
    FROM
        backend_auth.routes
	`

	rows, err := db.Pool().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routes, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Route])
	if err != nil {
		return nil, err
	}

	return routes, nil
}

func (db *DB) ExecuteManyRoutes(ctx context.Context, values []models.Route) error {
	query := `
	INSERT INTO
        backend_auth.routes (
            path,
            method
        )
    VALUES (
        $1,
        $2
    )
	`
	batch := &pgx.Batch{}
	for _, v := range values {
		batch.Queue(query, v.Path, v.Method)
	}

	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	br := tx.SendBatch(ctx, batch)
	for range values {
		_, err := br.Exec()
		if err != nil {
			_ = br.Close()
			return err
		}
	}
	err = br.Close()
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)

	return err
}
