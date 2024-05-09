package rest

import (
	"github.com/google/uuid"
	"practice_vgpek/internal/model/domain"
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
