package app

import "time"

type UserRole int

const (
	RoleGuest UserRole = 0
	RoleUser  UserRole = 1
	RoleAdmin UserRole = 2
)

func (role UserRole) String() string {
	switch role {
	case RoleGuest:
		return "guest"
	case RoleUser:
		return "user"
	case RoleAdmin:
		return "admin"
	default:
		return ""
	}
}

type User struct {
	Id         int
	Name       string
	Role       UserRole
	Registered time.Time
}
