package account

import (
	"context"
	"errors"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/entity"
)

type KeyDAO interface {
	ById(ctx context.Context, id int) (entity.Key, error)
}
type AccountDAO interface {
	ById(ctx context.Context, id int) (entity.Account, error)
}
type RoleDAO interface {
	ById(ctx context.Context, id int) (entity.Role, error)
}
type PermissionDAO interface {
	ByRoleId(ctx context.Context, roleId int) ([]entity.Permissions, error)
}

type Mediator struct {
	AccountDAO AccountDAO
	KeyDAO     KeyDAO
	RoleDAO    RoleDAO
	PermDAO    PermissionDAO
}

func NewAccountMediator(AccountDAO AccountDAO, KeyDAO KeyDAO,
	RoleDAO RoleDAO, PermDAO PermissionDAO) Mediator {
	return Mediator{
		AccountDAO: AccountDAO,
		KeyDAO:     KeyDAO,
		RoleDAO:    RoleDAO,
		PermDAO:    PermDAO,
	}
}

func (m Mediator) RoleByAccountId(ctx context.Context, id int) (entity.Role, error) {
	acc, err := m.AccountDAO.ById(ctx, id)
	if err != nil {
		return entity.Role{}, err
	}

	key, err := m.KeyDAO.ById(ctx, acc.KeyId)
	if err != nil {
		return entity.Role{}, err
	}

	role, err := m.RoleDAO.ById(ctx, key.RoleId)
	if err != nil {
		return entity.Role{}, err
	}

	return role, nil
}

func (m Mediator) HasAccess(ctx context.Context, accountId int, objectName, actionName string) (bool, error) {
	var hasObject, hasAction bool

	role, err := m.RoleByAccountId(ctx, accountId)
	if err != nil {
		return false, err
	}

	perms, err := m.PermDAO.ByRoleId(ctx, role.Id)
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
	acc, err := m.AccountDAO.ById(ctx, id)
	if err != nil {
		return domain.RolePermission{}, err
	}

	role, err := m.RoleDAO.ById(ctx, acc.RoleId)
	if err != nil {
		return domain.RolePermission{}, err
	}

	perms, err := m.PermDAO.ByRoleId(ctx, role.Id)
	if err != nil {
		return domain.RolePermission{}, err
	}

	var rolePerm domain.RolePermission

	var isDeleted bool

	if role.IsDeleted != nil {
		isDeleted = true
	}

	domainRole := domain.Role{
		ID:          role.Id,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		IsDeleted:   isDeleted,
		DeletedAt:   role.IsDeleted,
	}

	rolePerm.Role = domainRole

	for _, perm := range perms {
		var isDeleted bool

		if perm.Action.IsDeleted != nil {
			isDeleted = true
		}

		domainAction := domain.Action{
			ID:          perm.Action.Id,
			Name:        perm.Action.Name,
			Description: perm.Action.Description,
			CreatedAt:   perm.Action.CreatedAt,
			IsDeleted:   isDeleted,
			DeletedAt:   perm.Action.IsDeleted,
		}

		domainObject := domain.Object{
			ID:          perm.Object.Id,
			Name:        perm.Object.Name,
			Description: perm.Object.Description,
			CreatedAt:   perm.Object.CreatedAt,
			IsDeleted:   isDeleted,
			DeletedAt:   perm.Object.IsDeleted,
		}

		rolePerm.Object = domain.ObjectWithActions{
			Object:  domainObject,
			Actions: append(rolePerm.Object.Actions, domainAction),
		}
	}

	return rolePerm, nil
}
