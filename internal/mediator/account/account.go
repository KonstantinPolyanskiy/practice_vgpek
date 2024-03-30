package account

import (
	"context"
	"errors"
	"practice_vgpek/internal/model/account"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/registration_key"
)

type RepositoryKey interface {
	RegKeyById(ctx context.Context, id int) (registration_key.Entity, error)
}
type RepositoryAccount interface {
	AccountById(ctx context.Context, id int) (account.Entity, error)
}
type RepositoryRole interface {
	RoleById(ctx context.Context, id int) (permissions.RoleEntity, error)
}
type RepositoryPermission interface {
	PermissionsByRoleId(ctx context.Context, roleId int) ([]permissions.PermissionEntity, error)
}

type Mediator struct {
	AccountRepo RepositoryAccount
	RegKeyRepo  RepositoryKey
	RoleRepo    RepositoryRole
	PermRepo    RepositoryPermission
}

func NewAccountMediator(AccountRepo RepositoryAccount, RegKeyRepo RepositoryKey,
	RoleRepo RepositoryRole, PermRepo RepositoryPermission) Mediator {
	return Mediator{
		AccountRepo: AccountRepo,
		RegKeyRepo:  RegKeyRepo,
		RoleRepo:    RoleRepo,
		PermRepo:    PermRepo,
	}
}

func (m Mediator) RoleByAccountId(ctx context.Context, id int) (permissions.RoleEntity, error) {
	acc, err := m.AccountRepo.AccountById(ctx, id)
	if err != nil {
		return permissions.RoleEntity{}, err
	}

	key, err := m.RegKeyRepo.RegKeyById(ctx, acc.RegKeyId)
	if err != nil {
		return permissions.RoleEntity{}, err
	}

	role, err := m.RoleRepo.RoleById(ctx, key.RoleId)
	if err != nil {
		return permissions.RoleEntity{}, err
	}

	return role, nil
}

func (m Mediator) HasAccess(ctx context.Context, roleId int, objectName, actionName string) (bool, error) {
	var hasObject, hasAction bool

	perms, err := m.PermRepo.PermissionsByRoleId(ctx, roleId)
	if err != nil {
		return false, err
	}

	if len(perms) == 0 {
		return false, errors.New("no result")
	}

	for _, perm := range perms {
		if perm.ObjectEntity.Name == objectName {
			hasObject = true
		}
		if perm.ActionEntity.Name == actionName {
			hasAction = true
		}
	}

	return hasAction && hasObject, nil
}
