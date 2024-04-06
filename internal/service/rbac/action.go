package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
)

type ActionRepository interface {
	SaveAction(ctx context.Context, savingAction permissions.ActionDTO) (permissions.ActionEntity, error)
	ActionById(ctx context.Context, id int) (permissions.ActionEntity, error)
	ActionsByParams(ctx context.Context, params params.Default) ([]permissions.ActionEntity, error)
}

type AddedActionResult struct {
	Action permissions.AddActionResp
	Error  error
}

type GetActionResult struct {
	Action permissions.ActionEntity
	Error  error
}

type GetActionsResult struct {
	Actions []permissions.ActionEntity
	Error   error
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

func (s RBACService) ActionById(ctx context.Context, req permissions.GetActionReq) (permissions.ActionEntity, error) {
	resCh := make(chan GetActionResult)

	l := s.l.With(
		zap.String("operation", GetActionOperation),
		zap.String("layer", "services"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("error check access", zap.Error(err))

			sendGetActionResult(resCh, permissions.ActionEntity{}, ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetActionResult(resCh, permissions.ActionEntity{}, ErrDontHavePermission.Error())
			return
		}

		action, err := s.ar.ActionById(ctx, req.Id)
		if err != nil {
			sendGetActionResult(resCh, permissions.ActionEntity{}, "ошибка получения действия")
			return
		}

		sendGetActionResult(resCh, action, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return permissions.ActionEntity{}, ctx.Err()
		case result := <-resCh:
			return result.Action, result.Error

		}
	}
}

func (s RBACService) ActionsByParams(ctx context.Context, params params.Default) ([]permissions.ActionEntity, error) {
	resCh := make(chan GetActionsResult)

	l := s.l.With(
		zap.String("operation", GetActionsOperation),
		zap.String("layer", "services"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, GetActionName)
		if err != nil {
			l.Warn("error check access", zap.Error(err))

			sendGetActionsResult(resCh, nil, ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			sendGetActionsResult(resCh, nil, ErrDontHavePermission.Error())
			return
		}

		actions, err := s.ar.ActionsByParams(ctx, params)
		if err != nil {
			sendGetActionsResult(resCh, nil, "ошибка получения действий")
			return
		}

		sendGetActionsResult(resCh, actions, "")
		return

	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			return result.Actions, result.Error

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
func sendGetActionResult(resCh chan GetActionResult, resp permissions.ActionEntity, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetActionResult{
		Action: resp,
		Error:  err,
	}
}
func sendGetActionsResult(resCh chan GetActionsResult, resp []permissions.ActionEntity, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetActionsResult{
		Actions: resp,
		Error:   err,
	}
}
