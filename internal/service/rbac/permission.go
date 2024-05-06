package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
)

type PermissionDAO interface {
	Save(ctx context.Context, roleId, objectId int, actionsId []int) error
}

type AddedPermissionResult struct {
	Error error
}

func (s RBACService) NewPermission(ctx context.Context, req dto.SetPermissionReq) error {
	resCh := make(chan AddedPermissionResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.AddPermissionOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		if len(req.ActionsId) == 0 {
			l.Warn("нет действий для добавления")

			sendAddPermissionResult(resCh, "нет действий для добавления")
			return
		}

		err := s.permDAO.Save(ctx, req.RoleId, req.ObjectId, req.ActionsId)
		if err != nil {
			sendAddPermissionResult(resCh, "Неизвестная ошибка добавления доступов")
			return
		}

		sendAddPermissionResult(resCh, "")
	}()

	for {
		select {
		case <-ctx.Done():
			ctx.Err().Error()
		case result := <-resCh:
			result.Error.Error()
		}
	}
}

func sendAddPermissionResult(resCh chan AddedPermissionResult, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- AddedPermissionResult{
		Error: err,
	}
}
