package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	postType = iota
	attachmentType
	sharedType
)

func Connect(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	return &DB{db}, err
}

type DB struct {
	*sql.DB
}
