package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/account"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/model/registration_key"
	ar "practice_vgpek/internal/repository/account"
	pr "practice_vgpek/internal/repository/person"
	"practice_vgpek/internal/repository/rbac"
	kr "practice_vgpek/internal/repository/reg_key"
)

type PersonRepo interface {
	SavePerson(ctx context.Context, savingPerson person.DTO, accountId int) (person.Entity, error)
}

type PermissionRepo interface {
	SavePermission(ctx context.Context, roleId, objectId int, actionsId []int) error
	PermissionsByRoleId(ctx context.Context, roleId int) ([]permissions.PermissionEntity, error)
}

type ActionRepo interface {
	SaveAction(ctx context.Context, savingAction permissions.ActionDTO) (permissions.ActionEntity, error)
	ActionByName(ctx context.Context, name string) (permissions.ActionEntity, error)
	ActionById(ctx context.Context, id int) (permissions.ActionEntity, error)
	ActionsByParams(ctx context.Context, params params.Default) ([]permissions.ActionEntity, error)
}

type AccountRepo interface {
	SaveAccount(ctx context.Context, savingAcc account.DTO) (account.Entity, error)
	AccountByLogin(ctx context.Context, login string) (account.Entity, error)
	AccountById(ctx context.Context, id int) (account.Entity, error)
}

type ObjectRepo interface {
	SaveObject(ctx context.Context, savingObject permissions.ObjectDTO) (permissions.ObjectEntity, error)
}

type RoleRepo interface {
	SaveRole(ctx context.Context, savingRole permissions.RoleDTO) (permissions.RoleEntity, error)
	RoleByName(ctx context.Context, name string) (permissions.RoleEntity, error)
	RoleById(ctx context.Context, id int) (permissions.RoleEntity, error)
}

type KeyRepo interface {
	SaveKey(ctx context.Context, key registration_key.DTO) (registration_key.Entity, error)
	RegKeyByBody(ctx context.Context, body string) (registration_key.Entity, error)
	KeysByParams(ctx context.Context, params params.Key) ([]registration_key.Entity, error)
	RegKeyById(ctx context.Context, id int) (registration_key.Entity, error)
	IncCountUsages(ctx context.Context, keyId int) error
	Invalidate(ctx context.Context, keyId int) error
}

type Repository struct {
	l  *zap.Logger
	db *pgxpool.Pool
	PersonRepo
	KeyRepo
	AccountRepo
	ActionRepo
	ObjectRepo
	RoleRepo
	PermissionRepo
}

func New(db *pgxpool.Pool, logger *zap.Logger) Repository {
	return Repository{
		PersonRepo:     pr.NewPersonRepo(db, logger),
		KeyRepo:        kr.NewKeyRepo(db, logger),
		AccountRepo:    ar.NewAccountRepo(db, logger),
		ActionRepo:     rbac.NewActionRepo(db, logger),
		ObjectRepo:     rbac.NewObjectRepo(db, logger),
		RoleRepo:       rbac.NewRoleRepo(db, logger),
		PermissionRepo: rbac.NewPermissionRepository(db, logger),
	}
}
