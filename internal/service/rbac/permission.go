package rbac

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/permissions"
)

const (
	AddPermissionOperation = "добавление права действия в системе"
)

type PermissionRepo interface {
	SavePermission(ctx context.Context, roleId, objectId int, actionsId []int)
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
		saved, err := s.pr.SavePermission()/
	}()
}
