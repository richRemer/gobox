package app

import "time"

type User struct {
	Id         int
	Name       string
	Role       UserRole
	Registered time.Time
}
