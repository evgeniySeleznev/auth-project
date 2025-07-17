package model

import (
	"database/sql"
	"time"
)

type Role int32

const (
	RoleUnspecified Role = 0
	RoleUser        Role = 1
	RoleAdmin       Role = 2
)

// String возвращает строковое представление роли
func (r Role) String() string {
	switch r {
	case RoleUser:
		return "USER"
	case RoleAdmin:
		return "ADMIN"
	default:
		return "ROLE_UNSPECIFIED"
	}
}

// ParseRole преобразует строку в Role
func ParseRole(role string) Role {
	switch role {
	case "USER":
		return RoleUser
	case "ADMIN":
		return RoleAdmin
	default:
		return RoleUnspecified
	}
}

type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	Role      Role         `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
