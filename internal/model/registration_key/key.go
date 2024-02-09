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

type DTO struct {
	RoleId         int
	Body           string
	MaxCountUsages int
}

type Entity struct {
	RegKeyId           int
	RoleId             int
	Body               string
	MaxCountUsages     int
	CurrentCountUsages int
	CreatedAt          time.Time
	IsValid            bool
	InvalidationTime   time.Time
}
