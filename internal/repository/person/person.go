package person

import (
	"context"
	"practice_vgpek/internal/model/person"
	"time"
)

type Repository struct {
}

func NewPersonRepo() Repository {
	return Repository{}
}

func (r Repository) SavePerson(ctx context.Context, dto person.DTO) (person.Entity, error) {
	return person.Entity{
		Personality: person.Personality{
			FirstName:  dto.FirstName,
			MiddleName: dto.MiddleName,
			LastName:   dto.LastName,
		},
		PasswordHash:   dto.PasswordHash,
		Login:          dto.Login,
		CreatedAt:      time.Now(),
		IsActive:       true,
		DeactivateTime: time.Time{},
	}, nil
}
