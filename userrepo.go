package main

import (
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed schema/gobox.sql
var query string

type userrepo struct {
	db *sql.DB
}

func openUsers(dsn string) (*userrepo, error) {
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	repo := &userrepo{
		db: db,
	}

	return repo, nil
}

func (users userrepo) close() {
	users.db.Close()
}

func (users userrepo) initSchema() error {
	_, err := users.db.Exec(query)
	return err
}
