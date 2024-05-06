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
	"time"
)

type ObjectDAO interface {
	ById(ctx context.Context, id int) (entity.Object, error)
	Save(ctx context.Context, role dto.NewRBACPart) (entity.Object, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Object, error)
}

type GetObjectResult struct {
	Object domain.Object
	Error  error
}

type GetObjectsResult struct {
	Objects []domain.Object
	Error   error
}

type AddedObjectResult struct {
	Object domain.Object
	Error  error
}

func (s RBACService) NewObject(ctx context.Context, req dto.NewRBACReq) (domain.Object, error) {
	resCh := make(chan AddedObjectResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.AddObjectOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		// Проверяем что объект вообще введен
		if req.Name == "" {
			l.Warn("Пустой добавляемый объект")

			sendAddObjectResult(resCh, domain.Object{}, "Пустой добавляемый объект")
			return
		}

		part := dto.NewRBACPart{
			Name:        req.Name,
			Description: req.Description,
			CreatedAt:   time.Now(),
		}

		added, err := s.objectDAO.Save(ctx, part)
		if err != nil {
			sendAddObjectResult(resCh, domain.Object{}, "Неизвестная ошибка сохранения объекта действия")
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

		sendAddObjectResult(resCh, object, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Object{}, ctx.Err()
		case result := <-resCh:
			return result.Object, result.Error
		}
	}
}

func (s RBACService) ObjectById(ctx context.Context, req dto.EntityId) (domain.Object, error) {
	resCh := make(chan GetObjectResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.GetObjectOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		objectEntity, err := s.objectDAO.ById(ctx, req.Id)
		if err != nil {
			sendGetObjectResult(resCh, domain.Object{}, "Неизвестная ошибка получения объекта")
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

		sendGetObjectResult(resCh, object, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Object{}, ctx.Err()
		case result := <-resCh:
			return result.Object, result.Error
		}
	}
}

func (s RBACService) ObjectsByParams(ctx context.Context, p params.State) ([]domain.Object, error) {
	resCh := make(chan GetObjectsResult)

	_ = s.l.With(
		zap.String(operation.Operation, operation.GetObjectsOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		objectsEntity, err := s.objectDAO.ByParams(ctx, p.Default)
		if err != nil {
			sendGetObjectsResult(resCh, nil, "ошибка получения объектов действий")
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

		sendGetObjectsResult(resCh, objects, "")
		return

	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			return result.Objects, result.Error
		}
	}
}

func sendAddObjectResult(resCh chan AddedObjectResult, resp domain.Object, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- AddedObjectResult{
		Object: resp,
		Error:  err,
	}
}
func sendGetObjectResult(resCh chan GetObjectResult, resp domain.Object, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetObjectResult{
		Object: resp,
		Error:  err,
	}
}
func sendGetObjectsResult(resCh chan GetObjectsResult, resp []domain.Object, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetObjectsResult{
		Objects: resp,
		Error:   err,
	}
}
