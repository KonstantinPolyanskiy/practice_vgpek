package person

import (
	"practice_vgpek/internal/model/account"
	"time"
)

type Personality struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegistrationReq struct {
	Personality

	Credentials

	RegistrationKey string `json:"registration_key"`
}

type RegisteredResp struct {
	Personality
	CreatedAt time.Time `json:"created_at"`
}

type Entity struct {
	Personality
	PasswordHash   string
	Login          string
	CreatedAt      time.Time
	IsActive       bool
	DeactivateTime time.Time
}

type DTO struct {
	Personality
	Account account.DTO
}
