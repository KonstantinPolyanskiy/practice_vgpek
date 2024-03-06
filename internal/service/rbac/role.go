package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/permissions"
)

const (
	AddRoleAction = "добавление роли"
)

type RoleRepository interface {
	SaveRole(ctx context.Context, savingRole permissions.RoleDTO) (permissions.RoleEntity, error)
}

type AddedRoleResult struct {
	Role  permissions.AddRoleResp
	Error error
}

func (s RBACService) NewRole(ctx context.Context, addingRole permissions.AddRoleReq) (permissions.AddRoleResp, error) {
	resCh := make(chan AddedRoleResult)

	l := s.l.With(
		zap.String("action", AddRoleAction),
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
