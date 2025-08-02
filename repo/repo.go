package repo

import (
	"database/sql"
	_ "embed"
	"local/gobox/app"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed gobox.sql
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
	if users.db == nil {
		log.Printf("db is nil")
	} else {
		users.db.Close()
	}
}

func (users UserRepo) InitSchema() error {
	_, err := users.db.Exec(schema)
	return err
}

func (users UserRepo) InitSchemaIfNeeded() error {
	var name string

	query := `
		select name from sqlite_master
		where type = 'table' and name = 'user'`

	row := users.db.QueryRow(query)

	switch err := row.Scan(&name); err {
	case sql.ErrNoRows:
		return users.InitSchema()
	case nil:
		return nil
	default:
		return err
	}
}

func (users UserRepo) FindByName(name string) (app.User, error) {
	var user app.User
	var registered int64

	query := `
		select id, registered_at from user
		where name = ?`

	row := users.db.QueryRow(query, name)
	err := row.Scan(&user.Id, &registered)

	user.Name = name
	user.Role = app.RoleUser
	user.Registered = time.Unix(registered, 0)

	return user, err
}

func (users UserRepo) FindByPublicKey(pem string) (app.User, error) {
	var user app.User
	var registered int64

	query := `
		select u.id, u.name, u.registered_at
		from user u
		inner join public_key k on u.id = k.user_id
		where k.pem = ?`

	row := users.db.QueryRow(query, pem)
	err := row.Scan(&user.Id, &user.Name, &registered)

	user.Role = app.RoleUser
	user.Registered = time.Unix(registered, 0)

	return user, err
}

func (users UserRepo) FindOrRegister(name string) (app.User, error) {
	user, err := users.FindByName(name)

	if err != nil {
		user, err = users.Register(name)
	}

	return user, err
}

func (users UserRepo) Register(name string) (app.User, error) {
	var user app.User
	var registered int64

	query := `
		insert into user (name) values (?);

		select id, registered_at from user
		where id = last_insert_rowid();`

	row := users.db.QueryRow(query, name)
	err := row.Scan(&user.Id, &registered)

	user.Name = name
	user.Registered = time.Unix(registered, 0)

	return user, err
}

func (users UserRepo) RegisterWithKey(name string, pem string) (app.User, error) {
	var user app.User
	var registered int64
	var err error

	query := `
		begin;

		insert into user (name) values (:name);

		create table temp.t as
		select id, registered_at from user
		where id = last_insert_rowid();

		with inserted (id) as (select id from temp.t limit 1)
		insert into public_key (user_id, pem)
		select id, :pem from inserted;

		select * from temp.t;

		commit;`

	paramName := sql.Named("name", name)
	paramPem := sql.Named("pem", pem)
	rows, err := users.db.Query(query, paramName, paramPem)

	if err == nil {
		for rows.NextResultSet() {
			if rows.Next() {
				err = rows.Scan(&user.Id, registered)
				break
			}
		}
	}

	if err != nil {
		users.db.Exec("rollback")
		return user, err
	} else {
		users.db.Exec("drop table temp.t")
	}

	user.Name = name
	user.Registered = time.Unix(registered, 0)
	user.Role = app.RoleUser

	return user, nil
}
