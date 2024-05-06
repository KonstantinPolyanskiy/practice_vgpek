package rest

import (
	"practice_vgpek/internal/model/domain"
	"time"
)

type Key struct {
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

type Keys struct {
	Keys []Key `json:"keys"`
}

func (k Keys) DomainToResponse(keys []domain.Key) Keys {
	for _, key := range keys {
		k.Keys = append(k.Keys, Key{}.DomainToResponse(key))
	}

	return k
}

type InvalidatedKey struct {
	Id     int `json:"id"`
	RoleId int `json:"role_id"`

	CreatedAt time.Time `json:"created_at"`

	IsValid          bool      `json:"is_valid"`
	InvalidationTime time.Time `json:"invalidation_time"`
}

func (k InvalidatedKey) DomainToResponse(key domain.InvalidatedKey) InvalidatedKey {
	return InvalidatedKey{
		Id:               key.Id,
		RoleId:           key.RoleId,
		CreatedAt:        key.CreatedAt,
		IsValid:          key.IsValid,
		InvalidationTime: time.Now(),
	}
}

func (k Key) DomainToResponse(key domain.Key) Key {
	return Key{
		Id:             key.Id,
		RoleId:         key.RoleId,
		RoleName:       key.RoleName,
		Body:           key.Body,
		MaxCountUsages: key.MaxCountUsages,
		CountUsages:    key.CountUsages,
		CreatedAt:      key.CreatedAt,
		Group:          key.Group,
		IsValid:        key.IsValid,
	}
}
