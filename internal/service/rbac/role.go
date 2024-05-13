package rbac

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"time"
)

type RoleDAO interface {
	Save(ctx context.Context, part dto.NewRBACPart) (entity.Role, error)
	ById(ctx context.Context, id int) (entity.Role, error)
	SoftDeleteById(ctx context.Context, id int, info dto.DeleteInfo) error
	ByParams(ctx context.Context, p params.Default) ([]entity.Role, error)
}

func (s RBACService) NewRole(ctx context.Context, req dto.NewRBACReq) (domain.Role, error) {
	resCh := make(chan partResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.AddRoleOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		if req.Name == "" {
			l.Warn("Пустая добавляемая роль")

			sendPartResult(resCh, domain.Role{}, "Пустая добавляемая роль")
			return
		}

		part := dto.NewRBACPart{
			Name:        req.Name,
			Description: req.Description,
		}

		added, err := s.roleDAO.Save(ctx, part)
		if err != nil {
			sendPartResult(resCh, domain.Role{}, "Неизвестная ошибка сохранения роли")
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

		sendPartResult(resCh, role, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Role{}, ctx.Err()
		case result := <-resCh:
			return domain.Role(result.part.Part()), result.error
		}
	}
}

func (s RBACService) DeleteRoleById(ctx context.Context, req dto.EntityId) (domain.Role, error) {
	resCh := make(chan partResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.SoftDeleteRoleOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		info := dto.DeleteInfo{DeleteTime: time.Now()}

		err := s.roleDAO.SoftDeleteById(ctx, req.Id, info)
		if err != nil {
			l.Warn("ошибка мягкого удаления роли",
				zap.Int("id роли", req.Id),
				zap.Time("время удаления", info.DeleteTime),
			)

			sendPartResult(resCh, domain.Role{}, "Неизвестная ошибка удаления роли")
			return
		}

		deletedRoleEntity, err := s.roleDAO.ById(ctx, req.Id)
		if err != nil {
			l.Warn("ошибка получения удаленной роли", zap.Int("id роли", req.Id))

			sendPartResult(resCh, domain.Role{}, "Ошибка удаления роли")
			return
		}

		var isDeleted bool

		if deletedRoleEntity.IsDeleted != nil {
			isDeleted = true
		}

		role := domain.Role{
			ID:          deletedRoleEntity.Id,
			Name:        deletedRoleEntity.Name,
			Description: deletedRoleEntity.Description,
			CreatedAt:   deletedRoleEntity.CreatedAt,
			IsDeleted:   isDeleted,
			DeletedAt:   deletedRoleEntity.IsDeleted,
		}

		sendPartResult(resCh, role, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Role{}, ctx.Err()
		case result := <-resCh:
			return domain.Role(result.part.Part()), result.error
		}
	}
}

func (s RBACService) RoleById(ctx context.Context, req dto.EntityId) (domain.Role, error) {
	resCh := make(chan partResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.GetRoleOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		roleEntity, err := s.roleDAO.ById(ctx, req.Id)
		if err != nil {
			sendPartResult(resCh, domain.Role{}, "Ошибка получения роли")
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

		sendPartResult(resCh, role, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Role{}, ctx.Err()
		case result := <-resCh:
			return domain.Role(result.part.Part()), result.error
		}
	}
}

func (s RBACService) RolesByParams(ctx context.Context, p params.State) ([]domain.Role, error) {
	resCh := make(chan partsResult)

	_ = s.l.With(
		zap.String(operation.Operation, operation.GetRoleOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		rolesEntity, err := s.roleDAO.ByParams(ctx, p.Default)
		if err != nil {
			sendPartsResult(resCh, []domain.Role{}, "Ошибка получения ролей")
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
			resp = append(resp, roles...)
		case params.Deleted:
			resp = append(resp, filterDeleted(roles)...)
		case params.NotDeleted:
			resp = append(resp, filterNotDeleted(roles)...)
		}

		sendPartsResult(resCh, resp, "")
		return

	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			resp := make([]domain.Role, 0, len(result.parts))

			for _, role := range result.parts {
				resp = append(resp, domain.Role(role.Part()))
			}

			return resp, result.error
		}
	}
}
