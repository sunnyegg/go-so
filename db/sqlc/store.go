package db

import (
	"github.com/jackc/pgx/v5"
)

type Store interface {
	Querier
}

type SQLStore struct {
	db *pgx.Conn
	*Queries
}

func NewStore(db *pgx.Conn) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
