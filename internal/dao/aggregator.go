package dao

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/dao/account"
	"practice_vgpek/internal/dao/action"
	"practice_vgpek/internal/dao/issued"
	"practice_vgpek/internal/dao/key"
	"practice_vgpek/internal/dao/object"
	"practice_vgpek/internal/dao/permission"
	"practice_vgpek/internal/dao/person"
	"practice_vgpek/internal/dao/role"
	"practice_vgpek/internal/dao/solved"
)

type Aggregator struct {
	ActionDAO ActionDAO
	ObjectDAO ObjectDAO
	RoleDAO   RoleDAO

	PermissionDAO PermissionDAO

	PersonDAO  PersonDAO
	AccountDAO AccountDAO

	KeyDAO KeyDAO

	IssuedDAO IssuedPracticeDAO
	SolvedDAO SolvedPracticeDAO
}

func New(db *pgxpool.Pool, logger *zap.Logger) Aggregator {
	return Aggregator{
		ActionDAO: action.New(db, logger),
		ObjectDAO: object.New(db, logger),
		RoleDAO:   role.New(db, logger),

		PermissionDAO: permission.New(db, logger),

		PersonDAO:  person.New(db, logger),
		AccountDAO: account.New(db, logger),

		KeyDAO: key.New(db, logger),

		IssuedDAO: issued.New(db, logger),
		SolvedDAO: solved.New(db, logger),
	}
}
