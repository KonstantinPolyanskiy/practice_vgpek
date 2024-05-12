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
)

type ObjectDAO interface {
	ById(ctx context.Context, id int) (entity.Object, error)
	Save(ctx context.Context, role dto.NewRBACPart) (entity.Object, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Object, error)
}

func (s RBACService) NewObject(ctx context.Context, req dto.NewRBACReq) (domain.Object, error) {
	resCh := make(chan partResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.AddObjectOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		// Проверяем что объект вообще введен
		if req.Name == "" {
			l.Warn("Пустой добавляемый объект")

			sendPartResult(resCh, domain.Object{}, "Пустой добавляемый объект")
			return
		}

		part := dto.NewRBACPart{
			Name:        req.Name,
			Description: req.Description,
		}

		added, err := s.objectDAO.Save(ctx, part)
		if err != nil {
			sendPartResult(resCh, domain.Object{}, "Неизвестная ошибка сохранения объекта действия")
			return
		}

		object := domain.Object{
			ID:          added.Id,
			Name:        added.Name,
			Description: added.Description,
			CreatedAt:   added.CreatedAt,
			IsDeleted:   false,
			DeletedAt:   nil,
		}

		sendPartResult(resCh, object, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Object{}, ctx.Err()
		case result := <-resCh:
			return domain.Object(result.part.Part()), result.error
		}
	}
}

func (s RBACService) ObjectById(ctx context.Context, req dto.EntityId) (domain.Object, error) {
	resCh := make(chan partResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.GetObjectOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		objectEntity, err := s.objectDAO.ById(ctx, req.Id)
		if err != nil {
			sendPartResult(resCh, domain.Object{}, "Неизвестная ошибка получения объекта")
			return
		}

		object := domain.Object{
			ID:          objectEntity.Id,
			Name:        objectEntity.Name,
			Description: objectEntity.Description,
			CreatedAt:   objectEntity.CreatedAt,
			IsDeleted:   false,
			DeletedAt:   nil,
		}

		l.Info("получение объекта по id",
			zap.Int("id объекта", object.ID),
			zap.Time("время создания", object.CreatedAt),
			zap.Bool("удалено", object.IsDeleted),
		)

		sendPartResult(resCh, object, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Object{}, ctx.Err()
		case result := <-resCh:
			return domain.Object(result.part.Part()), result.error
		}
	}
}

func (s RBACService) ObjectsByParams(ctx context.Context, p params.State) ([]domain.Object, error) {
	resCh := make(chan partsResult)

	_ = s.l.With(
		zap.String(operation.Operation, operation.GetObjectsOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		objectsEntity, err := s.objectDAO.ByParams(ctx, p.Default)
		if err != nil {
			sendPartsResult(resCh, []domain.Object{}, "ошибка получения объектов действий")
			return
		}

		objects := make([]domain.Object, 0, len(objectsEntity))
		for _, objectEntity := range objectsEntity {
			var isDeleted bool

			if objectEntity.IsDeleted != nil {
				isDeleted = true
			}

			object := domain.Object{
				ID:          objectEntity.Id,
				Name:        objectEntity.Name,
				Description: objectEntity.Description,
				CreatedAt:   objectEntity.CreatedAt,
				IsDeleted:   isDeleted,
				DeletedAt:   objectEntity.IsDeleted,
			}

			objects = append(objects, object)
		}

		resp := make([]domain.Object, 0, len(objects))

		switch p.State {
		case params.All:
			copy(resp, objects)
		case params.Deleted:
			resp = append(resp, filterDeleted(objects)...)
		case params.NotDeleted:
			resp = append(resp, filterNotDeleted(objects)...)
		}

		sendPartsResult(resCh, objects, "")
		return

	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			resp := make([]domain.Object, 0, len(result.parts))

			for _, object := range result.parts {
				resp = append(resp, domain.Object(object.Part()))
			}

			return resp, result.error
		}
	}
}
