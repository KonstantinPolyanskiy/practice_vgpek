package entity

import "time"

type Role struct {
	Id int `db:"internal_role_id"`

	Description string `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	IsDeleted *time.Time `db:"is_deleted"`
}
