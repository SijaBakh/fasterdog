package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func New(dsn string, ctx context.Context) (*DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
