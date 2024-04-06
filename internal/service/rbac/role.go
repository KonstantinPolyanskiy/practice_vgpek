package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
)

type RoleRepository interface {
	SaveRole(ctx context.Context, savingRole permissions.RoleDTO) (permissions.RoleEntity, error)
	RoleById(ctx context.Context, id int) (permissions.RoleEntity, error)
	RolesByParams(ctx context.Context, params params.Default) ([]permissions.RoleEntity, error)
}

type AddedRoleResult struct {
	Role  permissions.AddRoleResp
	Error error
}

type GetRoleResult struct {
	Role  permissions.RoleEntity
	Error error
}

type GetRolesResult struct {
	Roles []permissions.RoleEntity
	Error error
}

func (s RBACService) NewRole(ctx context.Context, addingRole permissions.AddRoleReq) (permissions.AddRoleResp, error) {
	resCh := make(chan AddedRoleResult)

	l := s.l.With(
		zap.String("операция", operation.AddRoleOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		// Проверяем что роль вообще введена
		if addingRole.Name == "" {
			l.Warn("пустая роль для добавления")

			sendAddRoleResult(resCh, permissions.AddRoleResp{}, "Пустая роль для добавления")
			return
		}

		dto := permissions.RoleDTO{
			Name: addingRole.Name,
		}

		added, err := s.rr.SaveRole(ctx, dto)
		if err != nil {
			sendAddRoleResult(resCh, permissions.AddRoleResp{}, "Неизвестная ошибка сохранения роли")
			return
		}

		resp := permissions.AddRoleResp{
			Name: added.Name,
		}

		sendAddRoleResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return permissions.AddRoleResp{}, ctx.Err()
		case result := <-resCh:
			return result.Role, result.Error
		}
	}
}

func (s RBACService) RoleById(ctx context.Context, id int) (permissions.RoleEntity, error) {
	resCh := make(chan GetRoleResult)

	l := s.l.With(
		zap.String("операция", operation.GetRoleOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("ошибка проверки доступа", zap.Error(err))

			sendGetRoleResult(resCh, permissions.RoleEntity{}, permissions.ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetRoleResult(resCh, permissions.RoleEntity{}, permissions.ErrDontHavePerm.Error())
			return
		}

		role, err := s.rr.RoleById(ctx, id)
		if err != nil {
			sendGetRoleResult(resCh, permissions.RoleEntity{}, "Ошибка получения роли")
			return
		}

		sendGetRoleResult(resCh, role, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return permissions.RoleEntity{}, ctx.Err()
		case result := <-resCh:
			return result.Role, result.Error
		}
	}
}

func (s RBACService) RolesByParams(ctx context.Context, params params.Default) ([]permissions.RoleEntity, error) {
	resCh := make(chan GetRolesResult)

	l := s.l.With(
		zap.String("операция", operation.GetRoleOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("ошибка проверки доступа", zap.Error(err))

			sendGetRolesResult(resCh, nil, permissions.ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetRolesResult(resCh, nil, permissions.ErrDontHavePerm.Error())
			return
		}

		roles, err := s.rr.RolesByParams(ctx, params)
		if err != nil {
			sendGetRolesResult(resCh, nil, "Ошибка получения ролей")
			return
		}

		sendGetRolesResult(resCh, roles, "")
		return

	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			return result.Roles, result.Error
		}
	}
}

func sendAddRoleResult(resCh chan AddedRoleResult, resp permissions.AddRoleResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- AddedRoleResult{
		Role:  resp,
		Error: err,
	}
}
func sendGetRoleResult(resCh chan GetRoleResult, resp permissions.RoleEntity, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetRoleResult{
		Role:  resp,
		Error: err,
	}
}
func sendGetRolesResult(resCh chan GetRolesResult, resp []permissions.RoleEntity, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetRolesResult{
		Roles: resp,
		Error: err,
	}
}
