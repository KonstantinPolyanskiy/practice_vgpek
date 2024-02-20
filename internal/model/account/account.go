package account

import (
	"time"
)

type DTO struct {
	Login        string
	PasswordHash string
	RoleId       int
	RegKeyId     int
}

type Entity struct {
	AccountId      int        `db:"account_id"`
	Login          string     `db:"login"`
	PasswordHash   string     `db:"password_hash"`
	CreatedAt      time.Time  `db:"created_at"`
	IsActive       bool       `db:"is_active"`
	DeactivateTime *time.Time `db:"deactivate_time"`
	RoleId         int        `db:"internal_role_id"`
	RegKeyId       int        `db:"reg_key_id"`
}
