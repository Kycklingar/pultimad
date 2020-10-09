package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func connect(connStr string) error {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}
