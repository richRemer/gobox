package app

import "time"

type UserRole int

const (
	RoleGuest UserRole = 0
	RoleUser  UserRole = 1
	RoleAdmin UserRole = 2
)

type User struct {
	Id         int
	Name       string
	Role       UserRole
	Registered time.Time
}
