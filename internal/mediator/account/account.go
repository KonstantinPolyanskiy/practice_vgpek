package account

import (
	"context"
	"errors"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
)

type KeyService interface {
	ById(ctx context.Context, id dto.EntityId) (domain.Key, error)
}
type AccountService interface {
	AccountById(ctx context.Context, id dto.EntityId) (domain.Account, error)
}

type RoleService interface {
	ById(ctx context.Context, id dto.EntityId) (domain.Role, error)
}
type PermissionService interface {
	ByRoleId(ctx context.Context, id dto.EntityId) ([]domain.Permissions, error)
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

	role, err := m.RoleService.ById(ctx, dto.EntityId{Id: key.RoleId})
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

func (m Mediator) PermByAccountId(ctx context.Context, id int) (domain.RolePermission, error) {
	acc, err := m.AccountService.AccountById(ctx, dto.EntityId{Id: id})
	if err != nil {
		return domain.RolePermission{}, err
	}

	role, err := m.RoleService.ById(ctx, dto.EntityId{Id: acc.RoleId})
	if err != nil {
		return domain.RolePermission{}, err
	}

	perms, err := m.PermService.ByRoleId(ctx, dto.EntityId{Id: role.ID})
	if err != nil {
		return domain.RolePermission{}, err
	}

	var rolePerm domain.RolePermission

	domainRole := domain.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		IsDeleted:   role.IsDeleted,
		DeletedAt:   role.DeletedAt,
	}

	rolePerm.Role = domainRole

	for _, perm := range perms {

		domainAction := domain.Action{
			ID:          perm.Action.ID,
			Name:        perm.Action.Name,
			Description: perm.Action.Description,
			CreatedAt:   perm.Action.CreatedAt,
			IsDeleted:   perm.Action.IsDeleted,
			DeletedAt:   perm.Action.DeletedAt,
		}

		domainObject := domain.Object{
			ID:          perm.Object.ID,
			Name:        perm.Object.Name,
			Description: perm.Object.Description,
			CreatedAt:   perm.Object.CreatedAt,
			IsDeleted:   perm.Object.IsDeleted,
			DeletedAt:   perm.Object.DeletedAt,
		}

		rolePerm.Object = domain.ObjectWithActions{
			Object:  domainObject,
			Actions: append(rolePerm.Object.Actions, domainAction),
		}
	}

	return rolePerm, nil
}
