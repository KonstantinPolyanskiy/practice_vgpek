package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/model/registration_key"
	pr "practice_vgpek/internal/repository/person"
	kr "practice_vgpek/internal/repository/reg_key"
)

type PersonRepo interface {
	SavePerson(ctx context.Context, dto person.DTO) (person.Entity, error)
}

type KeyRepo interface {
	SaveKey(ctx context.Context, key registration_key.DTO) (registration_key.Entity, error)
	RegKeyByBody(ctx context.Context, body string) (registration_key.Entity, error)
}

type Repository struct {
	db *pgxpool.Pool
	PersonRepo
	KeyRepo
}

func New(db *pgxpool.Pool) Repository {
	return Repository{
		PersonRepo: pr.NewPersonRepo(db),
		KeyRepo:    kr.NewKeyRepo(db),
	}
}
