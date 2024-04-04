package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
)

const (
	AddObjectOperation  = "добавление объекта действия в системе"
	GetObjectOperation  = "получение объекта действия"
	GetObjectsOperation = "получение объектов действий"
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
		zap.String("action", AddObjectOperation),
		zap.String("layer", "services"),
	)

	go func() {
		// Проверяем что объект вообще введен
		if addingObject.Name == "" {
			l.Warn("empty adding object")
			sendAddObjectResult(resCh, permissions.AddObjectResp{}, "пустое имя объекта")
			return
		}

		dto := permissions.ObjectDTO{
			Name: addingObject.Name,
		}

		added, err := s.or.SaveObject(ctx, dto)
		if err != nil {
			l.Warn("error save object in db", zap.String("object name", addingObject.Name))
			sendAddObjectResult(resCh, permissions.AddObjectResp{}, "неизвестная ошибка сохранения объекта действия")
			return
		}

		l.Info("object successfully save", zap.String("object name", added.Name))

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
		zap.String("operation", GetObjectOperation),
		zap.String("layer", "services"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("error check access", zap.Error(err))

			sendGetObjectResult(resCh, permissions.ObjectEntity{}, ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetObjectResult(resCh, permissions.ObjectEntity{}, ErrDontHavePermission.Error())
			return
		}

		object, err := s.or.ObjectById(ctx, id)
		if err != nil {
			sendGetObjectResult(resCh, permissions.ObjectEntity{}, "ошибка получения объекта")
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
		zap.String("operation", GetObjectsOperation),
		zap.String("layer", "services"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("error check access", zap.Error(err))

			sendGetObjectsResult(resCh, nil, ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetObjectsResult(resCh, nil, ErrDontHavePermission.Error())
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
