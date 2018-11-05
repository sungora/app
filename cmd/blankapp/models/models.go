package models

import (
	"PKGAPPNAME/models/users"
)

func NewUser(id uint64) *users.User {
	return users.NewUser(id)
}
