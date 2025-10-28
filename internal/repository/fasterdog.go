package repository

import (
	"fasterdog/internal/adapter/db"
)

type FasterdogRepository struct {
	db *db.DB
}

func New(db *db.DB) *FasterdogRepository {
	return &FasterdogRepository{
		db: db,
	}
}
