package domain

import (
	"time"
)

type InvalidatedKey struct {
	Id     int
	RoleId int

	CreatedAt time.Time

	IsValid          bool
	InvalidationTime time.Time
}

type Key struct {
	Id int

	RoleId   int
	RoleName string

	Body string

	MaxCountUsages int
	CountUsages    int

	CreatedAt time.Time

	Group string

	IsValid bool
}
