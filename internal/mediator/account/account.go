package account

import (
	"context"
	"errors"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
)

type KeyService interface {
	ById(ctx context.Context, req dto.EntityId) (domain.Key, error)
}
type AccountService interface {
	AccountById(ctx context.Context, req dto.EntityId) (domain.Account, error)
}

type RoleService interface {
	RoleById(ctx context.Context, req dto.EntityId) (domain.Role, error)
}
type PermissionService interface {
	ByRoleId(ctx context.Context, req dto.EntityId) ([]domain.Permissions, error)
}

type Mediator struct {
	AccountService AccountService
	KeyService     KeyService
	RoleService    RoleService
	PermService    PermissionService
}

func NewAccountMediator(AccountService AccountService, KeyService KeyService,
	RoleService RoleService, PermService PermissionService) Mediator {
	return Mediator{
		AccountService: AccountService,
		KeyService:     KeyService,
		RoleService:    RoleService,
		PermService:    PermService,
	}
}

func (m Mediator) RoleByAccountId(ctx context.Context, id int) (domain.Role, error) {
	acc, err := m.AccountService.AccountById(ctx, dto.EntityId{Id: id})
	if err != nil {
		return domain.Role{}, err
	}

	key, err := m.KeyService.ById(ctx, dto.EntityId{Id: acc.KeyId})
	if err != nil {
		return domain.Role{}, err
	}

	role, err := m.RoleService.RoleById(ctx, dto.EntityId{Id: key.RoleId})
	if err != nil {
		return domain.Role{}, err
	}

	return role, nil
}

func (m Mediator) HasAccess(ctx context.Context, accountId int, objectName, actionName string) (bool, error) {
	var hasObject, hasAction bool

	role, err := m.RoleByAccountId(ctx, accountId)
	if err != nil {
		return false, err
	}

	perms, err := m.PermService.ByRoleId(ctx, dto.EntityId{Id: role.ID})
	if err != nil {
		return false, err
	}

	if len(perms) == 0 {
		return false, errors.New("no result")
	}

	for _, perm := range perms {
		if perm.Object.Name == objectName {
			hasObject = true
		}
		if perm.Action.Name == actionName {
			hasAction = true
		}
	}

	return hasAction && hasObject, nil
}
