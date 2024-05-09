package person

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
)

type PermResult struct {
	Perm  domain.RolePermission
	Error error
}

func (s Service) PermByAccountId(ctx context.Context, req dto.EntityId) (domain.RolePermission, error) {
	resCh := make(chan PermResult)

	_ = s.logger.With(
		zap.String(operation.Operation, operation.GetPermByAccountIdOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		role, err := s.roleService.RoleById(ctx, req)
		if err != nil {
			sendGetPermResult(resCh, domain.RolePermission{}, "ошибка получения роли")
			return
		}

		perms, err := s.permDAO.ByRoleId(ctx, role.ID)
		if err != nil {
			sendGetPermResult(resCh, domain.RolePermission{}, "ошибка получения доступов")
			return
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
			var actionIsDeleted, objectIsDeleted bool

			if perm.Action.IsDeleted != nil {
				actionIsDeleted = true
			}
			if perm.Action.IsDeleted != nil {
				objectIsDeleted = true
			}

			domainAction := domain.Action{
				ID:          perm.Action.Id,
				Name:        perm.Action.Name,
				Description: perm.Action.Description,
				CreatedAt:   perm.Action.CreatedAt,
				IsDeleted:   actionIsDeleted,
				DeletedAt:   perm.Action.IsDeleted,
			}

			domainObject := domain.Object{
				ID:          perm.Object.Id,
				Name:        perm.Object.Name,
				Description: perm.Object.Description,
				CreatedAt:   perm.Object.CreatedAt,
				IsDeleted:   objectIsDeleted,
				DeletedAt:   perm.Object.IsDeleted,
			}

			rolePerm.Object = domain.ObjectWithActions{
				Object:  domainObject,
				Actions: append(rolePerm.Object.Actions, domainAction),
			}
		}

		sendGetPermResult(resCh, rolePerm, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.RolePermission{}, ctx.Err()
		case result := <-resCh:
			return result.Perm, result.Error
		}
	}
}

func sendGetPermResult(ch chan PermResult, perm domain.RolePermission, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	ch <- PermResult{
		Perm:  perm,
		Error: err,
	}
}
