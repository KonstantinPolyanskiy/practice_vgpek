package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/permissions"
)

const (
	AddObjectAction = "добавление объекта действия в системе"
)

type ObjectRepository interface {
	SaveObject(ctx context.Context, savingObject permissions.ObjectDTO) (permissions.ObjectEntity, error)
}

type AddedObjectResult struct {
	Object permissions.AddObjectResp
	Error  error
}

func (s RBACService) NewObject(ctx context.Context, addingObject permissions.AddObjectReq) (permissions.AddObjectResp, error) {
	resCh := make(chan AddedObjectResult)

	l := s.l.With(
		zap.String("action", AddObjectAction),
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
