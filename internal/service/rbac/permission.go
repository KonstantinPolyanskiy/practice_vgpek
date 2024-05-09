package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
)

type PermissionDAO interface {
	Save(ctx context.Context, roleId, objectId int, actionsId []int) error
	ByRoleId(ctx context.Context, roleId int) ([]entity.Permissions, error)
}

type AddedPermissionResult struct {
	Error error
}

type GetPermByRoleIdResult struct {
	Permissions []domain.Permissions
	Error       error
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

func (s RBACService) ByRoleId(ctx context.Context, req dto.EntityId) ([]domain.Permissions, error) {
	resCh := make(chan GetPermByRoleIdResult)

	_ = s.l.With(
		zap.String(operation.Operation, operation.GetPermissionOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		permEntity, err := s.permDAO.ByRoleId(ctx, req.Id)
		if err != nil {
			sendGetPermByRoleIdResult(resCh, nil, "ошибка получения доступов")
			return
		}

		perm := make([]domain.Permissions, 0, len(permEntity))

		for _, p := range permEntity {
			var roleIsDeleted, objectIsDeleted, actionIsDeleted bool

			if p.Role.IsDeleted != nil {
				roleIsDeleted = true
			}
			if p.Object.IsDeleted != nil {
				objectIsDeleted = true
			}
			if p.Action.IsDeleted != nil {
				actionIsDeleted = true
			}

			perm = append(perm, domain.Permissions{
				PermissionId: p.PermissionId,
				Role: domain.Role{
					ID:          p.Role.Id,
					Name:        p.Role.Name,
					Description: p.Role.Description,
					CreatedAt:   p.Role.CreatedAt,
					IsDeleted:   roleIsDeleted,
					DeletedAt:   p.Role.IsDeleted,
				},
				Action: domain.Action{
					ID:          p.Action.Id,
					Name:        p.Action.Name,
					Description: p.Action.Description,
					CreatedAt:   p.Action.CreatedAt,
					IsDeleted:   actionIsDeleted,
					DeletedAt:   p.Action.IsDeleted,
				},
				Object: domain.Object{
					ID:          p.Object.Id,
					Name:        p.Object.Name,
					Description: p.Object.Description,
					CreatedAt:   p.Object.CreatedAt,
					IsDeleted:   objectIsDeleted,
					DeletedAt:   p.Object.IsDeleted,
				},
			})
		}

		sendGetPermByRoleIdResult(resCh, perm, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			return result.Permissions, result.Error
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

func sendGetPermByRoleIdResult(resCh chan GetPermByRoleIdResult, resp []domain.Permissions, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetPermByRoleIdResult{
		Permissions: resp,
		Error:       err,
	}
}
