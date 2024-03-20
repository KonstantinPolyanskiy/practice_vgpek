package person

import (
	"github.com/google/uuid"
	"practice_vgpek/internal/model/account"
	"time"
)

type Personality struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

type LogInReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LogInResp struct {
	Token string `json:"token"`
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
	PersonUUID uuid.UUID `db:"person_uuid"`
	FirstName  string    `db:"first_name"`
	MiddleName string    `db:"middle_name"`
	LastName   string    `db:"last_name"`
	AccountId  int       `db:"account_id"`
}

type DTO struct {
	Personality
	Account account.DTO
}
