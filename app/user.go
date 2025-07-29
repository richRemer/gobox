package app

import "time"

type User struct {
	Id         int
	Name       string
	Role       UserRole
	Registered time.Time
}

type UserRepo interface {
	Close()
	InitSchema() error
	InitSchemaIfNeeded() error
	FindByName(name string) (User, error)
	FindByPublicKey(pem string) (User, error)
	FindOrRegister(name string) (User, error)
	Register(name string) (User, error)
	RegisterWithKey(name string, pem string) (User, error)
}
