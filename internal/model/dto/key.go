package dto

import "time"

type NewKeyInfo struct {
	RoleId int

	Body           string
	MaxCountUsages int

	CreatedAt time.Time

	Group string
}

// NewKeyReq описывает данные, которые вводятся пользователем при создании нового ключа
type NewKeyReq struct {
	RoleId         int    `json:"role_id"`
	MaxCountUsages int    `json:"max_count_usages"`
	GroupName      string `json:"group_name"`
}

type KeyResp struct {
	Id             int       `json:"id"`
	RoleId         int       `json:"role_id"`
	RoleName       string    `json:"role_name"`
	Body           string    `json:"body"`
	MaxCountUsages int       `json:"max_count_usages"`
	CountUsages    int       `json:"count_usages"`
	CreatedAt      time.Time `json:"created_at"`
	Group          string    `json:"group"`
	IsValid        bool      `json:"is_valid"`
}

type DeleteKeyResp struct {
	Id     int `json:"id"`
	RoleId int `json:"role_id"`

	CreatedAt time.Time `json:"created_at"`

	IsValid          bool      `json:"is_valid"`
	InvalidationTime time.Time `json:"invalidation_time"`
}
