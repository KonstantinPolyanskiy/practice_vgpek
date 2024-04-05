package reg_key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/registration_key"
)

const (
	GetActionName = "GET"
)

type GetKeysResult struct {
	Keys  []registration_key.Entity
	Error error
}

func (s Service) KeysByParams(ctx context.Context, keyParams params.Key) ([]registration_key.Entity, error) {
	resCh := make(chan GetKeysResult)

	l := s.l.With(
		zap.String("операция", GetKeysOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, AddActionName)
		if err != nil || !hasAccess {
			l.Warn("возникла ошибка при проверке прав", zap.Error(err))

			sendGetKeysResult(resCh, nil, permissions.ErrDontHavePerm.Error())
			return
		}

		keys, err := s.r.KeysByParams(ctx, keyParams)
		if err != nil {
			sendGetKeysResult(resCh, nil, "Ошибка получения ключей")
			return
		}

		sendGetKeysResult(resCh, keys, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			return result.Keys, result.Error

		}
	}
}

func sendGetKeysResult(resCh chan GetKeysResult, resp []registration_key.Entity, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetKeysResult{
		Keys:  resp,
		Error: err,
	}
}
