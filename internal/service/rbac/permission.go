package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/permissions"
)

type PermissionRepo interface {
	SavePermission(ctx context.Context, roleId, objectId int, actionsId []int) error
}

type AddedPermissionResult struct {
	Perm  permissions.AddPermResp
	Error error
}

func (s RBACService) NewPermission(ctx context.Context, addingPerm permissions.AddPermReq) (permissions.AddPermResp, error) {
	resCh := make(chan AddedPermissionResult)

	l := s.l.With(
		zap.String("action", AddPermissionOperation),
		zap.String("layer", "services"),
	)

	go func() {
		if len(addingPerm.ActionsId) == 0 {
			l.Warn("пустая роль для добавления")
			sendAddPermissionResult(resCh, permissions.AddPermResp{}, "нет действий для добавления")
			return
		}

		err := s.pr.SavePermission(ctx, addingPerm.RoleId, addingPerm.ObjectId, addingPerm.ActionsId)
		if err != nil {
			sendAddPermissionResult(resCh, permissions.AddPermResp{}, "ошибка добавления доступов")
			return
		}

		sendAddPermissionResult(resCh, permissions.AddPermResp{}, "")
	}()

	for {
		select {
		case <-ctx.Done():
			return permissions.AddPermResp{}, ctx.Err()
		case result := <-resCh:
			return result.Perm, result.Error
		}
	}
}

func sendAddPermissionResult(resCh chan AddedPermissionResult, resp permissions.AddPermResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- AddedPermissionResult{
		Perm:  resp,
		Error: err,
	}
}
