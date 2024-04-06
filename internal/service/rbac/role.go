package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
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
		zap.String("action", AddRoleOperation),
		zap.String("layer", "services"),
	)

	go func() {
		// Проверяем что роль вообще введена
		if addingRole.Name == "" {
			l.Warn("empty adding role")
			sendAddRoleResult(resCh, permissions.AddRoleResp{}, "пустое название роли")
			return
		}

		dto := permissions.RoleDTO{
			Name: addingRole.Name,
		}

		added, err := s.rr.SaveRole(ctx, dto)
		if err != nil {
			l.Warn("error save role in db", zap.String("role name", addingRole.Name))
			sendAddRoleResult(resCh, permissions.AddRoleResp{}, "неизвестная ошибка сохранения роли")
			return
		}

		l.Info("role successfully save", zap.String("role name", added.Name))

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
		zap.String("operation", GetRoleOperation),
		zap.String("layer", "services"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("error check access", zap.Error(err))

			sendGetRoleResult(resCh, permissions.RoleEntity{}, ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetRoleResult(resCh, permissions.RoleEntity{}, ErrDontHavePermission.Error())
			return
		}

		role, err := s.rr.RoleById(ctx, id)
		if err != nil {
			sendGetRoleResult(resCh, permissions.RoleEntity{}, "ошибка получения роли")
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
		zap.String("operation", GetRoleOperation),
		zap.String("layer", "services"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("error check access", zap.Error(err))

			sendGetRolesResult(resCh, nil, ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetRolesResult(resCh, nil, ErrDontHavePermission.Error())
			return
		}

		roles, err := s.rr.RolesByParams(ctx, params)
		if err != nil {
			sendGetRolesResult(resCh, nil, "ошибка получения ролей")
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
