// internal/db/store.go
package db

import (
	"database/sql"
)

type SingleDb interface {
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) SingleDb {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}
