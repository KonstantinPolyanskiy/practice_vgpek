package rest

import (
	"github.com/google/uuid"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/entity"
	"time"
)

// Account описывает часть ответа сервиса на регистрацию пользователя
type Account struct {
	Login string `json:"login"`

	IsActive       bool       `json:"is_active"`
	DeactivateTime *time.Time `json:"deactivate_time"`

	RoleName string `json:"role_name"`
	RoleId   int    `json:"role_id"`

	CreatedAt time.Time `json:"created_at"`
}

type AccountEntity struct {
	Id             int        `json:"id"`
	Login          string     `json:"login"`
	PasswordHash   string     `json:"password_hash"`
	CreatedAt      time.Time  `json:"created_at"`
	IsActive       bool       `json:"is_active"`
	DeactivateTime *time.Time `json:"deactivate_time"`
	KeyId          int        `json:"key_id"`
	RoleId         int        `json:"role_id"`
}

func (a AccountEntity) EntityToResponse(account entity.Account) AccountEntity {
	return AccountEntity{
		Id:             account.Id,
		Login:          account.Login,
		PasswordHash:   account.PasswordHash,
		CreatedAt:      account.CreatedAt,
		IsActive:       account.IsActive,
		DeactivateTime: account.DeactivateTime,
		KeyId:          account.KeyId,
		RoleId:         account.RoleId,
	}
}

type AccountsEntity struct {
	Accounts []AccountEntity
}

func (a AccountsEntity) EntityToResponse(accounts []entity.Account) AccountsEntity {
	for _, acc := range accounts {
		a.Accounts = append(a.Accounts, AccountEntity{}.EntityToResponse(acc))
	}

	return a
}

type Person struct {
	Uuid uuid.UUID `json:"uuid"`

	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`

	Account `json:"account"`
}

func (p Person) DomainToResponse(user domain.Person) Person {
	return Person{
		Uuid:       user.UUID,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Account: Account{
			Login:          user.Login,
			IsActive:       user.IsActive,
			DeactivateTime: user.DeactivateTime,
			RoleName:       user.RoleName,
			RoleId:         user.RoleId,
			CreatedAt:      user.CreatedAt,
		},
	}
}
