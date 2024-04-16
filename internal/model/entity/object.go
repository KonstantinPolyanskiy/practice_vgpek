package entity

import "time"

type Object struct {
	Id int `db:"internal_object_id"`

	Name        string `db:"internal_object_name"`
	Description string `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	IsDeleted *time.Time `db:"is_deleted"`
}
