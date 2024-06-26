package rbac

import (
	"context"
	"go.uber.org/zap"
	"log"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"time"
)

type ActionDAO interface {
	Save(ctx context.Context, part dto.NewRBACPart) (entity.Action, error)
	ById(ctx context.Context, id int) (entity.Action, error)
	SoftDeleteById(ctx context.Context, id int, info dto.DeleteInfo) error
	ByParams(ctx context.Context, params params.Default) ([]entity.Action, error)
}

func (s RBACService) NewAction(ctx context.Context, req dto.NewRBACReq) (domain.Action, error) {
	resCh := make(chan partResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.AddActionOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		// Проверяем что действие - не пустая строка
		if req.Name == "" {
			l.Warn("попытка добавить пустое действие")

			sendPartResult(resCh, domain.Action{}, "Пустое добавляемое действие")
			return
		}

		// Формируем DTO
		part := dto.NewRBACPart{
			Name:        req.Name,
			Description: req.Description,
		}

		// Сохраняем действие в БД
		added, err := s.actionDAO.Save(ctx, part)
		if err != nil {
			sendPartResult(resCh, domain.Action{}, "Неизвестная ошибка сохранения действия")
			return
		}

		var isDeleted bool

		if added.IsDeleted != nil {
			isDeleted = true
		}

		// Формируем ответ
		action := domain.Action{
			ID:          added.Id,
			Name:        added.Name,
			Description: added.Description,
			CreatedAt:   added.CreatedAt,
			IsDeleted:   isDeleted,
			DeletedAt:   added.IsDeleted,
		}

		// Возвращаем ответ
		sendPartResult(resCh, action, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Action{}, ctx.Err()
		case result := <-resCh:
			return domain.Action(result.part.Part()), result.error
		}
	}
}

func (s RBACService) DeleteActionById(ctx context.Context, req dto.EntityId) (domain.Action, error) {
	resCh := make(chan partResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.SoftDeleteActionById),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		info := dto.DeleteInfo{
			DeleteTime: time.Now(),
		}

		err := s.actionDAO.SoftDeleteById(ctx, req.Id, info)
		if err != nil {
			l.Warn("возникла ошибка мягкого удаления действия",
				zap.Int("id", req.Id),
				zap.Time("время удаления", info.DeleteTime),
			)

			sendPartResult(resCh, domain.Action{}, "возникла ошибка удаления")
			return
		}

		deletedActionEntity, err := s.actionDAO.ById(ctx, req.Id)
		if err != nil {
			l.Warn("возникла ошибка получения удаленного действия", zap.Int("id", req.Id))

			sendPartResult(resCh, domain.Action{}, "возникла ошибка удаления")
			return
		}

		var isDeleted bool

		if deletedActionEntity.IsDeleted != nil {
			isDeleted = true
		}

		action := domain.Action{
			ID:          deletedActionEntity.Id,
			Name:        deletedActionEntity.Name,
			Description: deletedActionEntity.Description,
			CreatedAt:   deletedActionEntity.CreatedAt,
			IsDeleted:   isDeleted,
			DeletedAt:   deletedActionEntity.IsDeleted,
		}

		sendPartResult(resCh, action, "")
		return

	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Action{}, ctx.Err()
		case result := <-resCh:
			return domain.Action(result.part.Part()), result.error

		}
	}
}

func (s RBACService) ActionById(ctx context.Context, req dto.EntityId) (domain.Action, error) {
	resCh := make(chan partResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.GetActionOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		actionEntity, err := s.actionDAO.ById(ctx, req.Id)
		if err != nil {
			sendPartResult(resCh, domain.Action{}, "Ошибка получения действия")
			return
		}
		var isDeleted bool

		if actionEntity.IsDeleted != nil {
			isDeleted = true
		}

		log.Println(actionEntity.IsDeleted)

		// Формируем ответ
		action := domain.Action{
			ID:          actionEntity.Id,
			Name:        actionEntity.Name,
			Description: actionEntity.Description,
			CreatedAt:   actionEntity.CreatedAt,
			IsDeleted:   isDeleted,
			DeletedAt:   actionEntity.IsDeleted,
		}

		l.Info("получение действия по id",
			zap.Int("id действия", action.ID),
			zap.Time("время создания", actionEntity.CreatedAt),
			zap.Bool("удалено", isDeleted),
		)

		sendPartResult(resCh, action, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Action{}, ctx.Err()
		case result := <-resCh:
			return domain.Action(result.part.Part()), result.error
		}
	}
}

func (s RBACService) ActionsByParams(ctx context.Context, p params.State) ([]domain.Action, error) {
	resCh := make(chan partsResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.GetActionsOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		// Получаем действия из БД
		actionsEntity, err := s.actionDAO.ByParams(ctx, p.Default)
		if err != nil {
			sendPartsResult(resCh, []domain.Action{}, "ошибка получения действий")
			return
		}

		// Создаем слайс всех действий (доменов)
		actions := make([]domain.Action, 0, len(actionsEntity))
		for _, actionEntity := range actionsEntity {
			var isDeleted bool

			if actionEntity.IsDeleted != nil {
				isDeleted = true
			}

			action := domain.Action{
				ID:          actionEntity.Id,
				Name:        actionEntity.Name,
				Description: actionEntity.Description,
				CreatedAt:   actionEntity.CreatedAt,
				IsDeleted:   isDeleted,
				DeletedAt:   actionEntity.IsDeleted,
			}

			actions = append(actions, action)
		}

		resp := make([]domain.Action, 0, len(actions))

		switch p.State {
		case params.All:
			resp = append(resp, actions...)
		case params.Deleted:
			resp = append(resp, filterDeleted(actions)...)
		case params.NotDeleted:
			resp = append(resp, filterNotDeleted(actions)...)
		}

		l.Info("действия отданы", zap.Int("кол-во", len(resp)))

		sendPartsResult(resCh, resp, "")
		return

	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			resp := make([]domain.Action, 0, len(result.parts))
			for _, part := range result.parts {
				resp = append(resp, domain.Action(part.Part()))
			}

			return resp, result.error
		}
	}
}
