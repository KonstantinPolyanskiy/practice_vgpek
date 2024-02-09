package person

import "time"

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
	Login        string
	PasswordHash string
}
