package registration_key

import "time"

type AddReq struct {
	MaxCountUsages int `json:"max_count_usages"`
	RoleId         int `json:"role_id"`
}

type AddResp struct {
	RegKeyId           int       `json:"reg_key_id"`
	MaxCountUsages     int       `json:"max_count_usages"`
	CurrentCountUsages int       `json:"current_count_usages"`
	Body               string    `json:"body"`
	CreatedAt          time.Time `json:"created_at"`
}

type GetKeysResp struct {
	Keys []Entity `json:"keys"`
}

type DeleteReq struct {
	KeyId int `json:"key_id"`
}

type DeleteResp struct {
	KeyId int `json:"key_id"`
}

type DTO struct {
	RoleId         int
	Body           string
	MaxCountUsages int
}

type Entity struct {
	RegKeyId           int        `db:"reg_key_id"`
	RoleId             int        `db:"internal_role_id"`
	Body               string     `db:"body_key"`
	MaxCountUsages     int        `db:"max_count_usages"`
	CurrentCountUsages int        `db:"current_count_usages"`
	CreatedAt          time.Time  `db:"created_at"`
	IsValid            bool       `db:"is_valid"`
	InvalidationTime   *time.Time `db:"invalidation_time"`
}
