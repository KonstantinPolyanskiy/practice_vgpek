package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"practice_vgpek/internal/model/account"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/model/registration_key"
	ar "practice_vgpek/internal/repository/account"
	pr "practice_vgpek/internal/repository/person"
	kr "practice_vgpek/internal/repository/reg_key"
)

type PersonRepo interface {
	SavePerson(ctx context.Context, savingPerson person.DTO, accountId int) (person.Entity, error)
}

type AccountRepo interface {
	SaveAccount(ctx context.Context, savingAcc account.DTO) (account.Entity, error)
}

type KeyRepo interface {
	SaveKey(ctx context.Context, key registration_key.DTO) (registration_key.Entity, error)
	RegKeyByBody(ctx context.Context, body string) (registration_key.Entity, error)
	IncCountUsages(ctx context.Context, keyId int) error
	Invalidate(ctx context.Context, keyId int) error
}

type Repository struct {
	db *pgxpool.Pool
	PersonRepo
	KeyRepo
	AccountRepo
}

func New(db *pgxpool.Pool) Repository {
	return Repository{
		PersonRepo:  pr.NewPersonRepo(db),
		KeyRepo:     kr.NewKeyRepo(db),
		AccountRepo: ar.NewAccountRepo(db),
	}
}
