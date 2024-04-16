package entity

import "time"

type Role struct {
	Id int `db:"internal_role_id"`

	Name        string `db:"role_name"`
	Description string `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	IsDeleted *time.Time `db:"is_deleted"`
}
