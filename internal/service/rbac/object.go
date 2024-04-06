package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
)

type ObjectRepository interface {
	SaveObject(ctx context.Context, savingObject permissions.ObjectDTO) (permissions.ObjectEntity, error)
	ObjectById(ctx context.Context, id int) (permissions.ObjectEntity, error)
	ObjectsByParams(ctx context.Context, params params.Default) ([]permissions.ObjectEntity, error)
}

type GetObjectResult struct {
	Object permissions.ObjectEntity
	Error  error
}

type GetObjectsResult struct {
	Objects []permissions.ObjectEntity
	Error   error
}

type AddedObjectResult struct {
	Object permissions.AddObjectResp
	Error  error
}

func (s RBACService) NewObject(ctx context.Context, addingObject permissions.AddObjectReq) (permissions.AddObjectResp, error) {
	resCh := make(chan AddedObjectResult)

	l := s.l.With(
		zap.String("оперция", operation.AddObjectOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		// Проверяем что объект вообще введен
		if addingObject.Name == "" {
			l.Warn("Пустой добавляемый объект")

			sendAddObjectResult(resCh, permissions.AddObjectResp{}, "Пустой добавляемый объект")
			return
		}

		dto := permissions.ObjectDTO{
			Name: addingObject.Name,
		}

		added, err := s.or.SaveObject(ctx, dto)
		if err != nil {
			sendAddObjectResult(resCh, permissions.AddObjectResp{}, "Неизвестная ошибка сохранения объекта действия")
			return
		}

		resp := permissions.AddObjectResp{
			Name: added.Name,
		}

		sendAddObjectResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return permissions.AddObjectResp{}, ctx.Err()
		case result := <-resCh:
			return result.Object, result.Error
		}
	}
}

func (s RBACService) ObjectById(ctx context.Context, id int) (permissions.ObjectEntity, error) {
	resCh := make(chan GetObjectResult)

	l := s.l.With(
		zap.String("операция", operation.GetObjectOperation),
		zap.String("layer", "services"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("Ошибка проверки доступа", zap.Error(err))

			sendGetObjectResult(resCh, permissions.ObjectEntity{}, permissions.ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetObjectResult(resCh, permissions.ObjectEntity{}, permissions.ErrDontHavePerm.Error())
			return
		}

		object, err := s.or.ObjectById(ctx, id)
		if err != nil {
			sendGetObjectResult(resCh, permissions.ObjectEntity{}, "Неизвестная ошибка получения объекта")
			return
		}

		sendGetObjectResult(resCh, object, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return permissions.ObjectEntity{}, ctx.Err()
		case result := <-resCh:
			return result.Object, result.Error
		}
	}
}

func (s RBACService) ObjectsByParams(ctx context.Context, params params.Default) ([]permissions.ObjectEntity, error) {
	resCh := make(chan GetObjectsResult)

	l := s.l.With(
		zap.String("операция", operation.GetObjectsOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("ошибка проверки доступа", zap.Error(err))

			sendGetObjectsResult(resCh, nil, permissions.ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetObjectsResult(resCh, nil, permissions.ErrDontHavePerm.Error())
			return
		}

		objects, err := s.or.ObjectsByParams(ctx, params)
		if err != nil {
			sendGetObjectsResult(resCh, nil, "ошибка получения объектов действий")
			return
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

func sendAddObjectResult(resCh chan AddedObjectResult, resp permissions.AddObjectResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- AddedObjectResult{
		Object: resp,
		Error:  err,
	}
}
func sendGetObjectResult(resCh chan GetObjectResult, resp permissions.ObjectEntity, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetObjectResult{
		Object: resp,
		Error:  err,
	}
}
func sendGetObjectsResult(resCh chan GetObjectsResult, resp []permissions.ObjectEntity, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetObjectsResult{
		Objects: resp,
		Error:   err,
	}
}
