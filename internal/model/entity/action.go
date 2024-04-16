package entity

import "time"

type Action struct {
	Id int `db:"internal_action_id"`

	Name        string `db:"internal_action_name"`
	Description string `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	IsDeleted *time.Time `db:"is_deleted"`
}
