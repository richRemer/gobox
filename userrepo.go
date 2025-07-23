package main

import (
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed schema/gobox.sql
var schema string

type userrepo struct {
	db *sql.DB
}

func openUsers(dsn string) (*userrepo, error) {
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	return &userrepo{db: db}, nil
}

func (users userrepo) close() {
	users.db.Close()
}

func (users userrepo) initSchema() error {
	_, err := users.db.Exec(schema)
	return err
}
