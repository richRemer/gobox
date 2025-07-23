package repo

import (
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed schema/gobox.sql
var schema string

type UserRepo struct {
	db *sql.DB
}

func OpenUsers(dsn string) (*UserRepo, error) {
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	return &UserRepo{db: db}, nil
}

func (users UserRepo) Close() {
	users.db.Close()
}

func (users UserRepo) InitSchema() error {
	_, err := users.db.Exec(schema)
	return err
}

func (users UserRepo) InitSchemaIfNeeded() error {
	query := `
		select name from sqlite_master
		where type = 'table' and name = 'user'`

	row := users.db.QueryRow(query)

	switch err := row.Scan(); err {
	case sql.ErrNoRows:
		return users.InitSchema()
	case nil:
		return nil
	default:
		return err
	}
}
