package registration_key

import "time"

type AddReq struct {
	MaxCountUsages int    `json:"max_count_usages"`
	RoleId         int    `json:"role_id"`
	GroupName      string `json:"group_name"`
}

type AddResp struct {
	RegKeyId           int       `json:"reg_key_id"`
	MaxCountUsages     int       `json:"max_count_usages"`
	CurrentCountUsages int       `json:"current_count_usages"`
	Body               string    `json:"body"`
	CreatedAt          time.Time `json:"created_at"`
}

type GetKeyResp struct {
	RegKeyId           int        `json:"reg_key_id"`
	RoleId             int        `json:"role_id"`
	Body               string     `json:"body"`
	GroupName          string     `json:"group_name"`
	MaxCountUsages     int        `json:"max_count_usages"`
	CurrentCountUsages int        `json:"current_count_usages"`
	CreatedAt          time.Time  `json:"created_at"`
	IsValid            bool       `json:"is_valid"`
	InvalidationTime   *time.Time `json:"invalidation_time"`
}

type GetKeysResp struct {
	Keys []GetKeyResp `json:"keys"`
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
	GroupName      string
}

type Entity struct {
	RegKeyId           int        `db:"reg_key_id"`
	RoleId             int        `db:"internal_role_id"`
	Body               string     `db:"body_key"`
	GroupName          string     `db:"group_name"`
	MaxCountUsages     int        `db:"max_count_usages"`
	CurrentCountUsages int        `db:"current_count_usages"`
	CreatedAt          time.Time  `db:"created_at"`
	IsValid            bool       `db:"is_valid"`
	InvalidationTime   *time.Time `db:"invalidation_time"`
}
