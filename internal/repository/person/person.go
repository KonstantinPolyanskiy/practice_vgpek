package person

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"practice_vgpek/internal/model/person"
	"time"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewPersonRepo(db *pgxpool.Pool) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) SavePerson(ctx context.Context, dto person.DTO) (person.Entity, error) {
	return person.Entity{
		Personality: person.Personality{
			FirstName:  dto.FirstName,
			MiddleName: dto.MiddleName,
			LastName:   dto.LastName,
		},
		PasswordHash:   dto.Account.PasswordHash,
		Login:          dto.Account.Login,
		CreatedAt:      time.Now(),
		IsActive:       true,
		DeactivateTime: time.Time{},
	}, nil
}
