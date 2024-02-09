package repository

import (
	"context"
	"practice_vgpek/internal/model/person"
	pr "practice_vgpek/internal/repository/person"
)

type PersonRepo interface {
	SavePerson(ctx context.Context, dto person.DTO) (person.Entity, error)
}

type Repository struct {
	PersonRepo
}

func New() Repository {
	return Repository{
		PersonRepo: pr.NewPersonRepo(),
	}
}
