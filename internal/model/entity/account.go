package entity

import "time"

type Account struct {
	Id int `db:"account_id"`

	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`

	CreatedAt time.Time `db:"created_at"`

	IsActive       bool       `db:"is_active"`
	DeactivateTime *time.Time `db:"deactivate_time"`

	KeyId  int `db:"reg_key_id"`
	RoleId int `db:"internal_role_id"`
}
