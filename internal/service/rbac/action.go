package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/permissions"
)

const (
	AddActionOperation = "добавление права действия в системе"
)

type ActionRepository interface {
	SaveAction(ctx context.Context, savingAction permissions.ActionDTO) (permissions.ActionEntity, error)
}

type AddedActionResult struct {
	Action permissions.AddActionResp
	Error  error
}

func (s RBACService) NewAction(ctx context.Context, addingAction permissions.AddActionReq) (permissions.AddActionResp, error) {
	resCh := make(chan AddedActionResult)

	l := s.l.With(
		zap.String("action", AddActionOperation),
		zap.String("layer", "services"),
	)

	go func() {
		// Проверяем что действие - не пустая строка
		if addingAction.Name == "" {
			l.Warn("empty adding action")
			sendAddActionResult(resCh, permissions.AddActionResp{}, "пустая добавляемая роль")
			return
		}

		// Формируем DTO
		dto := permissions.ActionDTO{
			Name: addingAction.Name,
		}

		// Сохраняем действие в БД
		added, err := s.ar.SaveAction(ctx, dto)
		if err != nil {
			l.Warn("error save action in db", zap.String("action name", dto.Name))
			sendAddActionResult(resCh, permissions.AddActionResp{}, "неизвестная ошибка сохранения действия")
			return
		}

		l.Info("action successfully save", zap.String("action name", added.Name))

		// Формируем ответ
		resp := permissions.AddActionResp{
			Name: added.Name,
		}

		// Возвращаем ответ
		sendAddActionResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return permissions.AddActionResp{}, ctx.Err()
		case result := <-resCh:
			return result.Action, result.Error
		}
	}
}

func sendAddActionResult(resCh chan AddedActionResult, resp permissions.AddActionResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- AddedActionResult{
		Action: resp,
		Error:  err,
	}
}
