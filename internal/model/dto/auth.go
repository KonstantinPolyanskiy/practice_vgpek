package dto

import (
	"github.com/google/uuid"
	"time"
)

// RegistrationReq описывает данные, которые вводятся пользователем при регистрации аккаунта
type RegistrationReq struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name,omitempty"`

	Login    string `json:"login"`
	Password string `json:"password"`

	BodyKey string `json:"registration_key"`
}

// RegistrationResp ответ при успешной регистрации
type RegistrationResp struct {
	UUID uuid.UUID `json:"uuid"`

	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`

	Login string `json:"login"`

	IsActive       bool       `json:"is_active"`
	DeactivateTime *time.Time `json:"deactivate_time"`

	RoleName string `json:"role_name"`
	RoleId   int    `json:"role_id"`

	KeyId int `json:"key_id"`

	CreatedAt time.Time `json:"created_at"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// PersonRegistrationData вспомогательная структура, передающаяся на DAO слой для создания записи в БД
type PersonRegistrationData struct {
	UUID uuid.UUID

	FirstName  string
	SecondName string
	LastName   string

	AccountId int
}

// AccountRegistrationData вспомогательная структура, передающаяся на DAO слой для создания записи в БД
type AccountRegistrationData struct {
	Login        string
	PasswordHash string

	CreatedAt time.Time

	RoleId int
	KeyId  int
}
