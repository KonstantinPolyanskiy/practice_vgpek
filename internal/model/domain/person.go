package domain

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	Login string

	IsActive       bool
	DeactivateTime *time.Time

	RoleName string
	RoleId   int

	KeyId int

	CreatedAt time.Time
}

type Person struct {
	UUID uuid.UUID

	FirstName, MiddleName, LastName string

	Account
}
