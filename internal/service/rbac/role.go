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
	"practice_vgpek/internal/model/params"
)

type RoleDAO interface {
	Save(ctx context.Context, part dto.NewRBACPart) (entity.Role, error)
	ById(ctx context.Context, id int) (entity.Role, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Role, error)
}

type AddedRoleResult struct {
	Role  domain.Role
	Error error
}

type GetRoleResult struct {
	Role  domain.Role
	Error error
}

type GetRolesResult struct {
	Roles []domain.Role
	Error error
}

func (s RBACService) NewRole(ctx context.Context, req dto.NewRBACReq) (domain.Role, error) {
	resCh := make(chan AddedRoleResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.AddRoleOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		if req.Name == "" {
			l.Warn("Пустая добавляемая роль")

			sendAddRoleResult(resCh, domain.Role{}, "Пустая добавляемая роль")
			return
		}

		part := dto.NewRBACPart{
			Name:        req.Name,
			Description: req.Description,
		}

		added, err := s.roleDAO.Save(ctx, part)
		if err != nil {
			sendAddRoleResult(resCh, domain.Role{}, "Неизвестная ошибка сохранения роли")
			return
		}

		var isDeleted bool

		if added.IsDeleted != nil {
			isDeleted = true
		}

		// Формируем ответ
		role := domain.Role{
			ID:          added.Id,
			Name:        added.Name,
			Description: added.Description,
			CreatedAt:   added.CreatedAt,
			IsDeleted:   isDeleted,
			DeletedAt:   added.IsDeleted,
		}

		sendAddRoleResult(resCh, role, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Role{}, ctx.Err()
		case result := <-resCh:
			return result.Role, result.Error
		}
	}
}

func (s RBACService) RoleById(ctx context.Context, req dto.EntityId) (domain.Role, error) {
	resCh := make(chan GetRoleResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.GetRoleOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		roleEntity, err := s.roleDAO.ById(ctx, req.Id)
		if err != nil {
			sendGetRoleResult(resCh, domain.Role{}, "Ошибка получения роли")
			return
		}

		var isDeleted bool

		if roleEntity.IsDeleted != nil {
			isDeleted = true
		}

		role := domain.Role{
			ID:          roleEntity.Id,
			Name:        roleEntity.Name,
			Description: roleEntity.Description,
			CreatedAt:   roleEntity.CreatedAt,
			IsDeleted:   isDeleted,
			DeletedAt:   roleEntity.IsDeleted,
		}

		l.Info("получение роли по id",
			zap.Int("id роли", role.ID),
			zap.Time("время создания", role.CreatedAt),
			zap.Bool("удалено", isDeleted),
		)

		sendGetRoleResult(resCh, role, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Role{}, ctx.Err()
		case result := <-resCh:
			return result.Role, result.Error
		}
	}
}

func (s RBACService) RolesByParams(ctx context.Context, p params.State) ([]domain.Role, error) {
	resCh := make(chan GetRolesResult)

	_ = s.l.With(
		zap.String(operation.Operation, operation.GetRoleOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		rolesEntity, err := s.roleDAO.ByParams(ctx, p.Default)
		if err != nil {
			sendGetRolesResult(resCh, nil, "Ошибка получения ролей")
			return
		}

		roles := make([]domain.Role, 0, len(rolesEntity))
		for _, roleEntity := range rolesEntity {
			var isDeleted bool

			if roleEntity.IsDeleted != nil {
				isDeleted = true
			}

			role := domain.Role{
				ID:          roleEntity.Id,
				Name:        roleEntity.Name,
				Description: roleEntity.Description,
				CreatedAt:   roleEntity.CreatedAt,
				IsDeleted:   isDeleted,
				DeletedAt:   roleEntity.IsDeleted,
			}

			roles = append(roles, role)
		}

		resp := make([]domain.Role, 0, len(rolesEntity))

		switch p.State {
		case params.All:
			copy(resp, roles)
		case params.Deleted:
			resp = append(resp, filterDeleted(roles)...)
		case params.NotDeleted:
			resp = append(resp, filterNotDeleted(roles)...)
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

func sendAddRoleResult(resCh chan AddedRoleResult, resp domain.Role, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- AddedRoleResult{
		Role:  resp,
		Error: err,
	}
}
func sendGetRoleResult(resCh chan GetRoleResult, resp domain.Role, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetRoleResult{
		Role:  resp,
		Error: err,
	}
}
func sendGetRolesResult(resCh chan GetRolesResult, resp []domain.Role, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetRolesResult{
		Roles: resp,
		Error: err,
	}
}
