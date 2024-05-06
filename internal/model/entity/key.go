package entity

import "time"

type Key struct {
	Id     int `db:"reg_key_id"`
	RoleId int `db:"internal_role_id"`

	Body string `db:"body_key"`

	MaxCountUsages     int `db:"max_count_usages"`
	CurrentCountUsages int `db:"current_count_usages"`

	CreatedAt time.Time `db:"created_at"`

	IsValid          bool       `db:"is_valid"`
	InvalidationTime *time.Time `db:"invalidation_time"`

	GroupName string `db:"group_name"`
}
